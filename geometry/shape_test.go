package shape

import (
	"testing"
)

func TestCreatePoint(t *testing.T) {
	cases := []struct {
		want [2]float64
	}{
		{[2]float64{1.0, 1}},
		{[2]float64{1.1, 1e1}},
		{[2]float64{.0, -.12}},
		{[2]float64{0.213e1, -0.001}},
	}

	for _, c := range cases {
		got := NewPoint(c.want[0], c.want[1])
		if got.x != c.want[0] || got.y != c.want[1] {
			t.Errorf("TestFail, expected {x: %g, y: %g}", c.want[0], c.want[1])
		}
	}
}

func TestPointCompare(t *testing.T) {
	cases := []struct {
		f, s *Point
		want bool
	}{
		{NewPoint(1.0, 1), NewPoint(1.0000, 1.0), true},
		{NewPoint(1.000, 0.), NewPoint(1.0000, 0.0000000000), true},
		{NewPoint(1, 2), NewPoint(2, 1), false},
		{NewPoint(0, 0), NewPoint(.0, 0.000000000001), true},
	}

	for _, c := range cases {
		got := c.f.Equal(c.s)
		if got != c.want {
			t.Errorf("TestFail, expected %b", c.want)
		}
	}
}

func TestContainsPoint(t *testing.T) {
	cases := []struct {
		box   *BoundaryBox
		point *Point
		want  bool
	}{
		{NewBoundaryBox(NewPoint(0, 0), NewPoint(1, 1)), NewPoint(0.0001, 0.0), true},
		{NewBoundaryBox(NewPoint(0, 0), NewPoint(1, 1)), NewPoint(0.00000000001, -0.000000000001), false},
		{NewBoundaryBox(NewPoint(0, 0), NewPoint(0, 0)), NewPoint(0.000, 0.0), true},
		{NewBoundaryBox(NewPoint(-250, -10), NewPoint(250, 10)), NewPoint(0.0, 0.0), true},
		{NewBoundaryBox(NewPoint(0.00001, -0.00001), NewPoint(.00002, .00002)), NewPoint(0.00001, 0.0), true},
	}

	for _, c := range cases {
		got := c.box.ContainsPoint(c.point)
		if got != c.want {
			t.Errorf("TestFail, expected %b", c.want)
		}
	}
}

func TestIntersect(t *testing.T) {
	cases := []struct {
		box   *BoundaryBox
		other *BoundaryBox
		want  bool
	}{
		{NewBoundaryBox(NewPoint(0, 0), NewPoint(1, 1)), NewBoundaryBox(NewPoint(0, 0), NewPoint(1, 1)), true},
		{NewBoundaryBox(NewPoint(0, 0), NewPoint(1, 1)), NewBoundaryBox(NewPoint(1, 1), NewPoint(2, 2)), true},
		{NewBoundaryBox(NewPoint(0.5, 0.5), NewPoint(1, 1)), NewBoundaryBox(NewPoint(0, 0), NewPoint(0.5, 0.5)), true},
	}

	for _, c := range cases {
		got := c.box.Intersect(c.other)
		if got != c.want {
			t.Errorf("TestFail, expected %b", c.want)
		}
	}
}

func BenchmarkContaints(b *testing.B) {
	point := NewPoint(0.3, 0.2)
	box := NewBoundaryBox(NewPoint(0, 0), NewPoint(1, 1))

	for i := 0; i < b.N; i++ {
		box.ContainsPoint(point)
	}
}
