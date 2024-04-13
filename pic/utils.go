package pic

import (
	"image"
	"image/color"
)

func Equal(a, b image.Image) bool {
	if a.Bounds() != b.Bounds() {
		return false
	}
	bounds := a.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if !ColorEqual(a.At(x, y), b.At(x, y)) {
				return false
			}
		}
	}
	return true
}

func ColorEqual(a, b color.Color) bool {
	ra, ga, ba, aa := a.RGBA()
	rb, gb, bb, ab := b.RGBA()
	return ra == rb && ga == gb && ba == bb && aa == ab
}
