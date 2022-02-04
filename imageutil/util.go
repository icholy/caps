package imageutil

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strings"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/icholy/caps/shape"
	"github.com/nfnt/resize"
)

// Read reads an image
func Read(name string) (image.Image, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// MustRead reads an image and panics if it fails
func MustRead(name string) image.Image {
	m, err := Read(name)
	if err != nil {
		panic(err)
	}
	return m
}

// AverageColor returns the averaged RGB color from the image
func AverageColor(i image.Image) color.Color {
	var sumR, sumG, sumB, count uint32
	bounds := i.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := i.At(x, y).RGBA()
			if a != 0 {
				sumR += r
				sumG += g
				sumB += b
				count++
			}
		}
	}
	if count == 0 {
		return color.NRGBA{}
	}
	return color.NRGBA{
		R: uint8((sumR / count) >> 8),
		G: uint8((sumG / count) >> 8),
		B: uint8((sumB / count) >> 8),
		A: 255,
	}
}

// Crop the src image using the provided rect
func Crop(src image.Image, r shape.Rect) image.Image {
	bounds := r.ImageRect().Intersect(src.Bounds())
	m := image.NewRGBA(r.ImageRect())
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			m.Set(x, y, src.At(x, y))
		}
	}
	return m
}

// Resize an image to the specified width and height
func Resize(m image.Image, w, h float64) image.Image {
	return resize.Resize(uint(w), uint(h), m, resize.Lanczos3)
}

// ColorDistance returns a number indicating the difference between two colors
// NOTE: this approach probably sucks
func ColorDistance(c0, c1 color.Color) float64 {
	r0, g0, b0, _ := c0.RGBA()
	r1, g1, b1, _ := c1.RGBA()
	r := uintDist(r0, r1)
	g := uintDist(g0, g1)
	b := uintDist(b0, b1)
	return math.Sqrt(float64(r*r) + float64(g*g) + float64(b*b))
}

// uintDist returns the absolute difference between two values
func uintDist(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}

func ProminentColor(m image.Image) color.Color {
	items, err := prominentcolor.Kmeans(m)
	if err != nil {
		if strings.Contains(err.Error(), "Failed, no non-alpha pixels found") {
			return color.RGBA{}
		}
		panic(err)
	}
	item := items[0]
	return color.NRGBA{
		R: uint8(item.Color.R),
		G: uint8(item.Color.G),
		B: uint8(item.Color.B),
		A: 255,
	}
}

func ProminentColors(m image.Image, k int) []color.Color {
	items, err := prominentcolor.KmeansWithAll(k, m,
		prominentcolor.ArgumentDefault,
		prominentcolor.DefaultSize,
		prominentcolor.GetDefaultMasks(),
	)
	if err != nil {
		panic(err)
	}
	cc := make([]color.Color, len(items))
	for i, item := range items {
		cc[i] = color.NRGBA{
			R: uint8(item.Color.R),
			G: uint8(item.Color.G),
			B: uint8(item.Color.B),
			A: 255,
		}
	}
	return cc
}

func CloneColors(colors []color.Color) []color.Color {
	clone := make([]color.Color, len(colors))
	copy(clone, colors)
	return clone
}
