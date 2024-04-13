package geo

import (
	"bytes"
	_ "embed"
	"image"

	"context"
	"image/color"
	"image/png"
	"testing"

	"gene.lol/cmf"
	"gene.lol/cmf/pic"
	"github.com/stretchr/testify/assert"
)

//go:embed x-plus-y.png
var xPlusYBytes []byte
var xPlusY image.Image

func init() {
	xPlusY = cmf.Must(png.Decode(bytes.NewReader(xPlusYBytes)))
}

func TestRender(t *testing.T) {
	ft := FunctionTile[float64]{
		rows: 1024,
		cols: 1024,
		f: func(x, y int) float64 {
			return float64(x * y)
		},
	}
	img := Render(ft, ColorMap[float64]{
		FunctionEntry[float64]{
			MatchFunc: func(_ float64) bool { return true },
			ColorFunc: func(val float64) color.Color {
				intVal := uint32(val)
				return color.NRGBA{
					R: uint8(intVal >> 16),
					G: uint8(intVal >> 8),
					B: uint8(intVal),
					A: 0xFF,
				}
			},
		},
	})
	assert.Truef(t, pic.Equal(img, xPlusY), "something changed about rendering in geo.Render")
}

func TestRenderParallel(t *testing.T) {
	ft := FunctionTile[float64]{
		rows: 1024,
		cols: 1024,
		f: func(x, y int) float64 {
			x2 := x
			y2 := y
			return float64(x2 * y2)
		},
	}
	img := RenderParallel(context.Background(), ft, ColorMap[float64]{
		FunctionEntry[float64]{
			MatchFunc: func(_ float64) bool { return true },
			ColorFunc: func(val float64) color.Color {
				intVal := uint32(val)
				return color.NRGBA{
					R: uint8(intVal >> 16),
					G: uint8(intVal >> 8),
					B: uint8(intVal),
					A: 0xFF,
				}
			},
		},
	})
	assert.Truef(t, pic.Equal(img, xPlusY), "something changed about rendering in geo.RenderParallel")
}

// completely unscientific, I left it at the defaults and didn't shut down anything else
// but vaguely useful as a baseline
//
// go test -bench=.
// goos: darwin
// goarch: arm64
// pkg: gene.lol/cmf/geo
// BenchmarkRender-10            	       3	 449906472 ns/op
// BenchmarkRenderParallel-10    	      15	  69314189 ns/op
//
//
// go test -bench=. -benchtime=100s
// goos: darwin
// goarch: arm64
// pkg: gene.lol/cmf/geo
// BenchmarkRender-10            	     264	 459365071 ns/op
// BenchmarkRenderParallel-10    	    1712	  71252714 ns/op

func BenchmarkRender(b *testing.B) {
	ft := FunctionTile[float64]{
		rows: 4096,
		cols: 4096,
		f: func(x, y int) float64 {
			x2 := x
			y2 := y
			return float64(x2 * y2)
		},
	}
	b.ResetTimer()
	for range b.N {
		_ = Render(ft, ColorMap[float64]{
			FunctionEntry[float64]{
				MatchFunc: func(_ float64) bool { return true },
				ColorFunc: func(val float64) color.Color {
					intVal := uint32(val)
					return color.NRGBA{
						R: uint8(intVal >> 16),
						G: uint8(intVal >> 8),
						B: uint8(intVal),
						A: 0xFF,
					}
				},
			},
		})
	}
}

func BenchmarkRenderParallel(b *testing.B) {
	ft := FunctionTile[float64]{
		rows: 4096,
		cols: 4096,
		f: func(x, y int) float64 {
			x2 := x
			y2 := y
			return float64(x2 * y2)
		},
	}
	b.ResetTimer()
	for range b.N {
		_ = RenderParallel(context.Background(), ft, ColorMap[float64]{
			FunctionEntry[float64]{
				MatchFunc: func(_ float64) bool { return true },
				ColorFunc: func(val float64) color.Color {
					intVal := uint32(val)
					return color.NRGBA{
						R: uint8(intVal >> 16),
						G: uint8(intVal >> 8),
						B: uint8(intVal),
						A: 0xFF,
					}
				},
			},
		})
	}
}
