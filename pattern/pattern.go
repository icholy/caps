package pattern

import (
	"math"

	"github.com/icholy/caps/shape"
)

type Pattern interface {
	Circles() []shape.Circle
}

// Line creates a sequence of circles each with the specified
// delta values from the previous. The sequence stops a circle
// is out of bounds
type Line struct {
	Start  shape.Point
	Radius float64
	DeltaX float64
	DeltaY float64
	Bounds shape.Rect
}

// Circles returns a list of the line circles
func (l *Line) Circles() []shape.Circle {
	var (
		cc []shape.Circle
		c  = shape.Circle{Point: l.Start, R: l.Radius}
	)
	for l.Bounds.Overlaps(c) {
		cc = append(cc, c)
		c = c.Add(l.DeltaX, l.DeltaY)
	}
	return cc
}

// Square creates a sequence of circles that form
// a square grid inside the specified bounds
type Square struct {
	TopLeft shape.Point
	Bounds  shape.Rect
	Radius  float64
	Spacing float64
}

// Circles returns a list of the square grid circles
func (s *Square) Circles() []shape.Circle {
	var (
		cc    []shape.Circle
		delta = (2 * s.Radius) + s.Spacing
	)
	for p := s.TopLeft; s.Bounds.Overlaps(p); p = p.Add(delta, 0) {
		line := Line{
			Start:  p,
			Bounds: s.Bounds,
			Radius: s.Radius,
			DeltaY: delta,
		}
		cc = append(cc, line.Circles()...)
	}
	return cc
}

// TriangularV creates a triangluar tiled grid
// of circles inside the specified bounds
type TriangularV struct {
	TopLeft shape.Point
	Bounds  shape.Rect
	Radius  float64
	Spacing float64
}

// Circles returns the list of the triangular grid circles
func (t *TriangularV) Circles() []shape.Circle {
	var (
		cc     []shape.Circle
		deltaX = t.DeltaX()
		deltaY = 2*t.Radius + t.Spacing
		column int
	)
	for p := t.TopLeft; t.Bounds.Overlaps(p); p = p.Add(deltaX, 0) {
		var offsetY float64
		if column%2 == 0 {
			offsetY = t.Radius
		}
		line := Line{
			Start:  p.Sub(0, offsetY),
			Bounds: t.Bounds,
			Radius: t.Radius,
			DeltaY: deltaY,
		}
		cc = append(cc, line.Circles()...)
		column++
	}
	return cc
}

func (t *TriangularV) DeltaX() float64 {
	r := t.Radius + t.Spacing/2
	return math.Sqrt(math.Pow(2*r, 2) - math.Pow(r, 2))
}

// TriangularH creates a triangluar tiled grid
// of circles inside the specified bounds
type TriangularH struct {
	TopLeft shape.Point
	Bounds  shape.Rect
	Radius  float64
	Spacing float64
}

// Circles returns the list of the triangular grid circles
func (t *TriangularH) Circles() []shape.Circle {
	var (
		cc     []shape.Circle
		deltaY = t.DeltaY()
		deltaX = 2*t.Radius + t.Spacing
		column int
	)
	for p := t.TopLeft; t.Bounds.Overlaps(p); p = p.Add(0, deltaY) {
		var offsetX float64
		if column%2 == 0 {
			offsetX = t.Radius
		}
		line := Line{
			Start:  p.Sub(offsetX, 0),
			Bounds: t.Bounds,
			Radius: t.Radius,
			DeltaX: deltaX,
		}
		cc = append(cc, line.Circles()...)
		column++
	}
	return cc
}

func (t *TriangularH) DeltaY() float64 {
	r := t.Radius + t.Spacing/2
	return math.Sqrt(math.Pow(2*r, 2) - math.Pow(r, 2))
}
