package geo

import (
	"context"
	"image"
	"image/color"
	"sync"

	"github.com/GeneTinderholm/cmf/config"
	cmfConstraints "github.com/GeneTinderholm/cmf/constraints"

	"golang.org/x/exp/constraints"
)

type Tile[T any] interface {
	Dim() (rows, cols int)
	Get(x, y int) T
}

func ToArrayTile[T any](t Tile[T]) ArrayTile[T] {
	if at, ok := t.(ArrayTile[T]); ok {
		return at
	}
	rows, cols := t.Dim()
	data := make([]T, rows*cols)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			data[y*cols+x] = t.Get(x, y)
		}
	}
	return ArrayTile[T]{rows: rows, cols: cols, data: data}
}
func NewArrayTile[T any](rows, cols int, data []T) Tile[T] {
	return ArrayTile[T]{rows: rows, cols: cols, data: data}
}

type ArrayTile[T any] struct {
	data       []T
	rows, cols int
}

func (at ArrayTile[T]) Dim() (rows, cols int) {
	return at.rows, at.cols
}
func (at ArrayTile[T]) Get(x, y int) T {
	return at.data[y*at.cols+x]
}
func (at ArrayTile[T]) Set(x, y int, val T) {
	at.data[y*at.cols+x] = val
}
func (at ArrayTile[T]) Map(f func(x, y int, val T) T) ArrayTile[T] {
	newTile := ArrayTile[T]{rows: at.rows, cols: at.cols, data: make([]T, len(at.data))}
	for i, val := range at.data {
		newTile.data[i] = f(i%at.cols, i/at.cols, val)
	}
	return newTile
}
func (at ArrayTile[T]) LazyMap(f func(x, y int, val T) T) Tile[T] {
	return FunctionTile[T]{
		rows: at.rows,
		cols: at.cols,
		f: func(x, y int) T {
			return f(x, y, at.data[y*at.cols+x])
		},
	}
}
func (at ArrayTile[T]) ToArray() []T {
	return at.data
}

func NewFunctionTile[T any](rows, cols int, f func(int, int) T) FunctionTile[T] {
	return FunctionTile[T]{f: f, rows: rows, cols: cols}
}

type FunctionTile[T any] struct {
	f          func(int, int) T
	rows, cols int
}

func (ft FunctionTile[T]) Dim() (rows, cols int) {
	return ft.rows, ft.cols
}
func (ft FunctionTile[T]) Get(x, y int) T {
	return ft.f(x, y)
}

type ColorMapEntry[T any] interface {
	Matches(T) bool
	Color(T) color.Color
}
type ColorMap[T any] []ColorMapEntry[T]

func Render[T any](t Tile[T], cm ColorMap[T]) image.Image {
	rows, cols := t.Dim()
	img := image.NewRGBA(image.Rect(0, 0, cols, rows))
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			renderCell(t, x, y, cm, img)
		}
	}
	return img
}

func RenderParallel[T any](ctx context.Context, t Tile[T], cm ColorMap[T]) image.Image {
	rows, cols := t.Dim()
	img := image.NewRGBA(image.Rect(0, 0, cols, rows))

	numWorkers := config.GetConfig(ctx).Parallelism
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := range numWorkers {
		go func() {
			defer wg.Done()
			for y := i; y < rows; y += numWorkers {
				for x := 0; x < cols; x++ {
					renderCell(t, x, y, cm, img)
				}
			}
		}()
	}
	wg.Wait()

	return img
}

func renderCell[T any](t Tile[T], x, y int, cm ColorMap[T], img *image.RGBA) {
	f := t.Get(x, y)
	for _, entry := range cm {
		if entry.Matches(f) {
			img.Set(x, y, entry.Color(f))
			break
		}
	}
}

type ScalarColorMapEntry[T comparable] struct {
	Val T
	C   color.Color
}

func (s ScalarColorMapEntry[T]) Matches(f T) bool {
	return s.Val == f
}
func (s ScalarColorMapEntry[T]) Color(_ float64) color.Color {
	return s.C
}

type ColorRampEntry[T constraints.Ordered] struct {
	Min, Max T
	C        color.Color
}

func (ramp ColorRampEntry[T]) Matches(f T) bool {
	return ramp.Min <= f && f <= ramp.Max
}
func (ramp ColorRampEntry[T]) Color(f float64) color.Color {
	return ramp.C
}

type GradientEntry[T cmfConstraints.Number] struct {
	Min, Max       T
	ColorA, ColorB color.Color
}

func (ge GradientEntry[T]) Matches(f T) bool {
	return ge.Min <= f && f <= ge.Max
}
func (ge GradientEntry[T]) Color(f T) color.Color {
	if f <= ge.Min {
		return ge.ColorA
	}
	if f >= ge.Max {
		return ge.ColorB
	}
	return scaleColor(ge.ColorA, ge.ColorB, float64((f-ge.Min)/(ge.Max/ge.Min)))
}
func scaleColor(a, b color.Color, percent float64) color.Color {
	ra, ga, ba, aa := a.RGBA()
	rb, gb, bb, ab := b.RGBA()
	return color.RGBA64{
		R: uint16(scaleChannel(ra, rb, percent)),
		G: uint16(scaleChannel(ga, gb, percent)),
		B: uint16(scaleChannel(ba, bb, percent)),
		A: uint16(scaleChannel(aa, ab, percent)),
	}
}
func scaleChannel(a, b uint32, percent float64) uint32 {
	af := float64(a)
	bf := float64(b)
	diff := bf - af
	return uint32(af + (percent * diff))
}

type FunctionEntry[T any] struct {
	MatchFunc func(T) bool
	ColorFunc func(T) color.Color
}

func (fe FunctionEntry[T]) Matches(f T) bool {
	return fe.MatchFunc(f)
}
func (fe FunctionEntry[T]) Color(f T) color.Color {
	return fe.ColorFunc(f)
}
