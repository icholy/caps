package shape

import (
	"fmt"
	"image"
)

// Shape is satisfied by all shape types
type Shape interface {
	Bounds() Rect
	String() string
}

// Point represents a the X, Y position in an image
// with 0, 0 being at the top left
type Point struct {
	X, Y float64
}

// NewPoint is constructs a Point
func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

// String returns the string representation of the point
func (p Point) String() string {
	return fmt.Sprintf("Point(X:%f, Y:%f)", p.X, p.Y)
}

// Bounds returns the points Bounds
func (p Point) Bounds() Rect {
	return Rect{Min: p, Max: p}
}

// Add adds the x, y values to the Point
func (p Point) Add(x, y float64) Point {
	return Point{
		X: p.X + x,
		Y: p.Y + y,
	}
}

// Sub subtracts the x, y values from the Point
func (p Point) Sub(x, y float64) Point {
	return Point{
		X: p.X - x,
		Y: p.Y - y,
	}
}

// Rect is a rectangle represented by it's min an max points
type Rect struct {
	Min Point
	Max Point
}

// NewRect constructs a new rectangle
func NewRect(xMin, yMin, xMax, yMax float64) Rect {
	return Rect{
		Min: NewPoint(xMin, yMin),
		Max: NewPoint(xMax, yMax),
	}
}

// ImageBounds returns the image.Rectangle as a Rect
func FromImageRect(r image.Rectangle) Rect {
	return Rect{
		Min: Point{
			X: float64(r.Min.X),
			Y: float64(r.Min.Y),
		},
		Max: Point{
			X: float64(r.Max.X),
			Y: float64(r.Max.Y),
		},
	}
}

// ImageRect returns the rectangle as an image.Rectangle
func (r Rect) ImageRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(r.Min.X),
			Y: int(r.Min.Y),
		},
		Max: image.Point{
			X: int(r.Max.X),
			Y: int(r.Max.Y),
		},
	}
}

// Grow returns r with all dimentions grown by w.
// Pass a negative value to shrink.
func (r Rect) Grow(w float64) Rect {
	return Rect{
		Min: Point{
			X: r.Min.X - w,
			Y: r.Min.Y - w,
		},
		Max: Point{
			X: r.Max.X + w,
			Y: r.Max.Y + w,
		},
	}
}

// String returns the string representation of the rect
func (r Rect) String() string {
	return fmt.Sprintf("Rect(Min:%s, Max:%s)", r.Min, r.Max)
}

// Bounds returns the rect itself
func (r Rect) Bounds() Rect {
	return r
}

// W returns the rectangles width
func (r Rect) W() float64 {
	return r.Max.X - r.Min.X
}

// H returns the rectangles height
func (r Rect) H() float64 {
	return r.Max.Y - r.Min.Y
}

// Contains returns true when the shape's bounds
// are contained within the rectangle
func (r Rect) Contains(s Shape) bool {
	b := s.Bounds()
	return r.Min.X <= b.Min.X &&
		r.Min.Y <= b.Min.Y &&
		r.Max.X >= b.Max.X &&
		r.Max.Y >= b.Max.Y
}

// Overlaps returns true when the shape's bounds
// overlap with the rectangle
func (r Rect) Overlaps(s Shape) bool {
	b := s.Bounds()
	if b.Max.X < r.Min.X || r.Max.X < b.Min.X {
		return false
	}
	if b.Max.Y < r.Min.Y || r.Max.Y < b.Min.Y {
		return false
	}
	return true
}

// Circle consists of a center point and a radius
type Circle struct {
	Point
	R float64
}

// NewCircle constructs a new circle
func NewCircle(x, y, r float64) Circle {
	return Circle{Point: NewPoint(x, y), R: r}
}

// String returns the string representation of the circle
func (c Circle) String() string {
	return fmt.Sprintf("Circle(X:%f, Y:%f, R:%f)", c.X, c.Y, c.R)
}

// Add the x, y values to the circles center point
func (c Circle) Add(x, y float64) Circle {
	return Circle{
		Point: c.Point.Add(x, y),
		R:     c.R,
	}
}

// Sub subtracts the x, y values form the circle center point
func (c Circle) Sub(x, y float64) Circle {
	return Circle{
		Point: c.Point.Sub(x, y),
		R:     c.R,
	}
}

// Bounds returns the circles bounds
func (c Circle) Bounds() Rect {
	return Rect{
		Min: c.Point.Sub(c.R, c.R),
		Max: c.Point.Add(c.R, c.R),
	}
}
