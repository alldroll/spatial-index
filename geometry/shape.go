package shape

import (
    "fmt"
    "time"
)

var _ = fmt.Printf
var _ = time.Sleep

type Shape interface {
}

/**/
type Point struct {
    x, y float64
}

/**/
func NewPoint(x, y float64) *Point{
    return &Point{x, y}
}

/**/
func (self *Point) EqualXY(x, y float64) bool {
    return self.x == x && self.y == y
}

/* */
func (self *Point) Equal(other *Point) bool {
    return self.EqualXY(other.x, other.y)
}

/**/
type BoundaryBox struct {
    bl, tr *Point /*bottom left, top right*/
    points  []*Point
}

func NewBoundaryBox(bl, tr *Point) *BoundaryBox {
    return &BoundaryBox{bl, tr, []*Point{}}
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
    return [4]*BoundaryBox{
        NewBoundaryBox(
            NewPoint(self.bl.x, self.bl.y + (self.tr.y - self.bl.y) / 2),
            NewPoint(self.bl.x + (self.tr.x - self.bl.x) / 2, self.tr.y),
        ),
        NewBoundaryBox(
            NewPoint(
                self.bl.x + (self.tr.x - self.bl.x) / 2,
                self.bl.y + (self.tr.y - self.bl.y) / 2,
            ),
            self.tr,
        ),
        NewBoundaryBox(
            self.bl,
            NewPoint(
                self.bl.x + (self.tr.x - self.bl.x) / 2,
                self.bl.y + (self.tr.y - self.bl.y) / 2,
            ),
        ),
        NewBoundaryBox(
            NewPoint(self.bl.x + (self.tr.x - self.bl.x) / 2, self.bl.y),
            NewPoint(self.tr.x, self.bl.y + (self.tr.y - self.bl.y) / 2),
        ),
    }
}

func (self *BoundaryBox) Area() float64 {
    return (self.tr.x - self.bl.x) * (self.tr.y - self.bl.y)
}

func (self *BoundaryBox) AppendPoint(point *Point) {
    self.points = append(self.points, point)
}

func (self *BoundaryBox) GetPoints() []*Point {
    fmt.Printf(
        "BOX x1: %f, y1: %f, x2: %f, y2: %f\n",
        self.bl.x, self.bl.y, self.tr.x, self.tr.y,
    )
    return self.points
}

func (self *BoundaryBox) Intersect(other *BoundaryBox) bool {
    return self.tr.x >= other.bl.x && other.tr.x >= self.bl.x &&
        self.tr.y >= other.bl.y && other.tr.y >= self.bl.y
}

func (self *BoundaryBox) GetPointsCount() int {
    return len(self.points)
}