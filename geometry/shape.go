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
	*BoundaryBox
	count int
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

func (self *Point) Copy() *Point {
	return NewPoint(self.x, self.y)
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

func (self *BoundaryBox) Copy() *BoundaryBox {
	return NewBoundaryBox(self.bl, self.tr)
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

func (self *BoundaryBox) Extend(other *BoundaryBox) *BoundaryBox {
	return self.extend(other.bl, other.tr)
}

func (self *BoundaryBox) ExtendPoint(point *Point) *BoundaryBox {
	return self.extend(point, point)
}

func (self *BoundaryBox) extend(bl *Point, tr *Point) *BoundaryBox {
	if self.ContainsPoint(bl) && self.ContainsPoint(tr) {
		return self
	}

	return NewBoundaryBox(
		NewPoint(math.Min(self.bl.x, bl.x), math.Min(self.bl.y, bl.y)),
		NewPoint(math.Max(self.tr.x, tr.x), math.Max(self.tr.y, tr.y)),
	)
}

func NewCluster(point *Point, count int) *Cluster {
	return &Cluster{NewBoundaryBox(point, point), count}
}

func (self *Cluster) Copy() *Cluster {
	return &Cluster{self.BoundaryBox, self.count}
}

func (self *Cluster) GetX() float64 {
	return (self.bl.x + self.tr.x) / 2
}

func (self *Cluster) GetY() float64 {
	return (self.bl.y + self.tr.y) / 2
}

func (self *Cluster) GetCenter() *Point {
	return NewPoint(self.GetX(), self.GetY())
}

func (self *Cluster) GetCount() int {
	return self.count
}

func (self *Cluster) AddPoint(point *Point) {
	self.count++
	self.extend(point, point)
}
