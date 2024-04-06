package geo

import (
	"image"
	"image/color"
)

type Tile interface {
	Dim() (rows, cols int)
	Get(x, y int) float64
}

func ToArrayTile(t Tile) ArrayTile {
	if at, ok := t.(ArrayTile); ok {
		return at
	}
	rows, cols := t.Dim()
	data := make([]float64, rows*cols)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			data[y*cols+x] = t.Get(x, y)
		}
	}
	return ArrayTile{rows: rows, cols: cols, data: data}
}
func NewArrayTile(rows, cols int, data []float64) Tile {
	return ArrayTile{rows: rows, cols: cols, data: data}
}

type ArrayTile struct {
	data       []float64
	rows, cols int
}

func (at ArrayTile) Dim() (rows, cols int) {
	return at.rows, at.cols
}

func (at ArrayTile) Get(x, y int) float64 {
	return at.data[y*at.cols+x]
}
func (at ArrayTile) ToArray() []float64 {
	return at.data
}

func NewFunctionTile(rows, cols int, f func(int, int) float64) FunctionTile {
	return FunctionTile{f: f, rows: rows, cols: cols}
}

type FunctionTile struct {
	f          func(int, int) float64
	rows, cols int
}

func (ft FunctionTile) Dim() (rows, cols int) {
	return ft.rows, ft.cols
}
func (ft FunctionTile) Get(x, y int) float64 {
	return ft.f(x, y)
}

type ColorMapEntry interface {
	Matches(float64) bool
	Color(float64) color.Color
}
type ColorMap []ColorMapEntry

func Render(t Tile, cm ColorMap) image.Image {
	rows, cols := t.Dim()
	img := image.NewRGBA(image.Rect(0, 0, cols, rows))
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			f := t.Get(x, y)
			for _, entry := range cm {
				if entry.Matches(f) {
					img.Set(x, y, entry.Color(f))
					break
				}
			}
		}
	}
	return img
}

type ScalarColorMapEntry struct {
	Val float64
	C   color.Color
}

func (s ScalarColorMapEntry) Matches(f float64) bool {
	return s.Val == f
}
func (s ScalarColorMapEntry) Color(_ float64) color.Color {
	return s.C
}

type ColorRampEntry struct {
	Min, Max float64
	C        color.Color
}

func (ramp ColorRampEntry) Matches(f float64) bool {
	return ramp.Min <= f && f <= ramp.Max
}
func (ramp ColorRampEntry) Color(f float64) color.Color {
	return ramp.C
}

type GradientEntry struct {
	Min, Max       float64
	ColorA, ColorB color.Color
}

func (ge GradientEntry) Matches(f float64) bool {
	return ge.Min <= f && f <= ge.Max
}
func (ge GradientEntry) Color(f float64) color.Color {
	if f <= ge.Min {
		return ge.ColorA
	}
	if f >= ge.Max {
		return ge.ColorB
	}
	return scaleColor(ge.ColorA, ge.ColorB, (f-ge.Min)/(ge.Max/ge.Min))
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
