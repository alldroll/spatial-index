package shape

import "math"

const tolerance = 0.000001

type Point struct {
	x, y float64
}

type BoundaryBox struct {
	bl, tr *Point /*bottom left, top right*/
}

type Cluster struct {
	center *Point
	count  int
}

func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

func (self *Point) EqualXY(x, y float64) bool {
	return math.Abs(self.x-x) < tolerance && math.Abs(self.y-y) < tolerance
}

func (self *Point) Equal(other *Point) bool {
	return self.EqualXY(other.x, other.y)
}

func (self *Point) GetX() float64 {
	return self.x
}

func (self *Point) GetY() float64 {
	return self.y
}

func (self *Point) Plus(other *Point) {
	self.x += other.x
	self.y += other.y
}

func NewBoundaryBox(bl, tr *Point) *BoundaryBox {
	return &BoundaryBox{bl, tr}
}

func (self *BoundaryBox) GetBottomLeft() *Point {
	return self.bl
}

func (self *BoundaryBox) GetTopRight() *Point {
	return self.tr
}

func (self *BoundaryBox) Equal(other *BoundaryBox) bool {
	return self.bl.Equal(other.bl) && self.tr.Equal(other.tr)
}

func (self *BoundaryBox) ContainsPoint(point *Point) bool {
	return self.bl.x <= point.x && self.bl.y <= point.y &&
		self.tr.x >= point.x && self.tr.y >= point.y
}

func (self *BoundaryBox) ContainsBox(boundary *BoundaryBox) bool {
	return self.ContainsPoint(boundary.bl) && self.ContainsPoint(boundary.tr)
}

func (self *BoundaryBox) Quarter() [4]*BoundaryBox {
	var xm, ym float64 = (self.tr.x - self.bl.x) / 2, (self.tr.y - self.bl.y) / 2

	return [4]*BoundaryBox{
		NewBoundaryBox(
			NewPoint(self.bl.x, self.bl.y+ym),
			NewPoint(self.bl.x+xm, self.tr.y),
		),
		NewBoundaryBox(
			NewPoint(self.bl.x+xm, self.bl.y+ym),
			self.tr,
		),
		NewBoundaryBox(
			self.bl,
			NewPoint(self.bl.x+xm, self.bl.y+ym),
		),
		NewBoundaryBox(
			NewPoint(self.bl.x+xm, self.bl.y),
			NewPoint(self.tr.x, self.bl.y+ym),
		),
	}
}

func (self *BoundaryBox) Area() float64 {
	return (self.tr.x - self.bl.x) * (self.tr.y - self.bl.y)
}

func (self *BoundaryBox) Intersect(other *BoundaryBox) bool {
	return self.tr.x >= other.bl.x && other.tr.x >= self.bl.x &&
		self.tr.y >= other.bl.y && other.tr.y >= self.bl.y
}

func NewCluster(point *Point, count int) *Cluster {
	return &Cluster{point, count}
}

func (self *Cluster) GetX() float64 {
	res := self.center.GetX()
	if self.count > 0 {
		res = res / float64(self.count)
	}

	return res
}

func (self *Cluster) GetY() float64 {
	res := self.center.GetY()
	if self.count > 0 {
		res = res / float64(self.count)
	}

	return res
}

func (self *Cluster) GetCenter() *Point {
	return NewPoint(self.GetX(), self.GetY())
}

func (self *Cluster) GetCount() int {
	return self.count
}

func (self *Cluster) AddPoint(point *Point) {
	self.count++
	self.center.Plus(point)
}
