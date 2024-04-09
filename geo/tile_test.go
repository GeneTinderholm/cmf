package geo

import (
	"image/color"
	"image/png"
	"os"
	"testing"

	"gene.lol/cmf"
)


func TestRender(t *testing.T) {
	ft := FunctionTile[float64]{
		rows: 256,
		cols:256,
		f: func(x, y int) float64 {
			x2 := x / 4
			y2 := y / 4
			return float64(x2 * y2 * x2 * y2)
		},
	}
	img := Render(ft, ColorMap[float64]{
		FunctionEntry[float64]{
			MatchFunc: func(_ float64) bool {return true},
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
	f := cmf.Must(os.Create("test.png"))
	cmf.CheckErr(png.Encode(f, img))
}
