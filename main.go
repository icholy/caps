package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"golang.org/x/image/colornames"

	"github.com/icholy/caps/bottlecap"
	"github.com/icholy/caps/imageutil"
	"github.com/icholy/caps/pattern"
	"github.com/icholy/caps/shape"
)

type Canvas struct {
	ctx *gg.Context
}

func NewCanvas(bounds shape.Rect) *Canvas {
	ctx := gg.NewContext(int(bounds.W()), int(bounds.H()))
	ctx.DrawRectangle(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y)
	ctx.SetColor(colornames.White)
	ctx.Fill()
	return &Canvas{
		ctx: ctx,
	}
}

func (cv *Canvas) DrawCircleColor(c shape.Circle, cl color.Color) {
	cv.ctx.DrawCircle(c.X, c.Y, c.R)
	cv.ctx.SetColor(cl)
	cv.ctx.Fill()
}

func (cv *Canvas) DrawCircle(c shape.Circle) {
	cv.DrawCircleColor(c, colornames.Black)
}

func (cv *Canvas) DrawImage(m image.Image) {
	cv.ctx.DrawImage(m, 0, 0)
}

func (cv *Canvas) DrawImageAt(m image.Image, p shape.Point) {
	cv.ctx.DrawImage(m, int(p.X), int(p.Y))
}

func (cv *Canvas) SavePNG(name string) error {
	return cv.ctx.SavePNG(name)
}

func main() {

	const radius = 20
	const drawCaps = true
	const padding = -5

	// read all the caps and resize them
	caps := bottlecap.MustReadDir("images/caps", radius)

	// read the source image
	src := imageutil.MustRead("images/rubik_cube.png")

	// find the K main colors in the image and assign them to bottle caps
	pallet := imageutil.ProminentColors(src, len(caps))
	bottlecap.Shuffle(caps)
	bottlecap.AssignBestColors(caps, pallet)

	// get the bounds of the image as a rect
	bounds := shape.FromImageRect(src.Bounds())

	// create a canvas to draw on
	cv := NewCanvas(bounds)

	// specify the bottle cap layout pattern
	p := pattern.TriangularV{
		TopLeft: shape.NewPoint(40, 40),
		Bounds:  bounds,
		Radius:  radius,
		Spacing: 0,
	}

	circles := p.Circles()

	fmt.Printf("Drawing %d caps\n", len(circles))

	for _, c := range p.Circles() {
		bounds := c.Bounds()

		// find the average color for the square where the
		// the circle will be drawn
		cropped := imageutil.Crop(src, bounds.Grow(padding))
		avg := imageutil.AverageColor(cropped)

		// don't draw on tiles with nothing in them
		if _, _, _, a := avg.RGBA(); a == 0 {
			continue
		}

		// find the best matching cap for that color
		cap := bottlecap.BestMatch(avg, caps)
		if drawCaps {
			cv.DrawImageAt(cap.Image, bounds.Min)
		} else {
			cv.DrawCircleColor(c, cap.Color)
		}
	}

	// write it out
	if err := cv.SavePNG("out.png"); err != nil {
		log.Fatal(err)
	}
}
