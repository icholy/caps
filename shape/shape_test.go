package shape

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPointAdd(t *testing.T) {
	require.EqualValues(t, NewPoint(1, 2), NewPoint(0, 0).Add(1, 2))
	require.EqualValues(t, NewPoint(10, 2), NewPoint(5, 1).Add(5, 1))
	require.EqualValues(t, NewPoint(-1, -2), NewPoint(0, 0).Sub(1, 2))
}

func TestContains(t *testing.T) {
	bounds := NewRect(0, 0, 100, 100)
	tests := []struct {
		shape    Shape
		expected bool
	}{
		// points
		{NewPoint(5, 5), true},
		{NewPoint(0, 0), true},
		{NewPoint(100, 100), true},
		{NewPoint(101, 50), false},
		{NewPoint(50, 101), false},
		{NewPoint(-1, 50), false},
		{NewPoint(50, -1), false},

		// rects
		{NewRect(10, 10, 90, 90), true},
		{NewRect(0, 0, 0, 0), true},
		{NewRect(100, 100, 100, 100), true},
		{NewRect(-1, 10, 90, 90), false},
		{NewRect(10, -1, 90, 90), false},
		{NewRect(10, 10, 101, 90), false},
		{NewRect(10, 10, 100, 101), false},

		// circles
		{NewCircle(50, 50, 10), true},
		{NewCircle(10, 10, 5), true},
		{NewCircle(0, 0, 1), false},
		{NewCircle(0, 100, 1), false},
		{NewCircle(100, 100, 1), false},
		{NewCircle(100, 0, 1), false},
	}

	for _, tt := range tests {
		require.Equal(t, tt.expected, bounds.Contains(tt.shape), tt.shape)
	}
}

func TestOverlaps(t *testing.T) {
	bounds := NewRect(0, 0, 100, 100)
	tests := []struct {
		shape    Shape
		expected bool
	}{
		// points
		{NewPoint(5, 5), true},
		{NewPoint(0, 0), true},
		{NewPoint(100, 100), true},
		{NewPoint(101, 50), false},
		{NewPoint(50, 101), false},
		{NewPoint(-1, 50), false},
		{NewPoint(50, -1), false},

		// rects
		{NewRect(10, 10, 90, 90), true},
		{NewRect(0, 0, 0, 0), true},
		{NewRect(100, 100, 100, 100), true},
		{NewRect(-1, 10, 90, 90), true},
		{NewRect(10, -1, 90, 90), true},
		{NewRect(10, 10, 101, 90), true},
		{NewRect(10, 10, 100, 101), true},
		{NewRect(-10, -10, -1, -1), false},
		{NewRect(101, 101, 110, 110), false},

		// circles
		{NewCircle(50, 50, 10), true},
		{NewCircle(10, 10, 5), true},
		{NewCircle(0, 0, 1), true},
		{NewCircle(0, 100, 1), true},
		{NewCircle(100, 100, 1), true},
		{NewCircle(100, 0, 1), true},
		{NewCircle(110, 0, 1), false},
		{NewCircle(50, -50, 49), false},
	}

	for _, tt := range tests {
		require.Equal(t, tt.expected, bounds.Overlaps(tt.shape), tt.shape)
	}
}
