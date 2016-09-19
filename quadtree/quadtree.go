package quadtree

import (
	//"fmt"
	"github.com/alldroll/quadtree/geometry"
	"log"
)

/**/
type Part int

const (
	TOP_LEFT     Part = 0
	TOP_RIGHT         = 1
	BOTTOM_LEFT       = 2
	BOTTOM_RIGHT      = 3
)

const CHILDREN_COUNT = 4

/**/
type Node struct {
	box      *shape.BoundaryBox
	points   []*shape.Point
	children [CHILDREN_COUNT]*Node
	length   int
	level    int
	capacity int
}

/**/
type QuadTree struct {
	root   *Node
	length int
}

/**/
type QuadTreeError struct {
	msg string
}

func (e *QuadTreeError) Error() string {
	return e.msg
}

func NewNode(box *shape.BoundaryBox, level int, capacity int) *Node {
	return &Node{
		box,
		[]*shape.Point{},
		[CHILDREN_COUNT]*Node{nil, nil, nil, nil},
		0,
		level,
		capacity,
	}
}

func (node *Node) SplitNode() {
	boxes := node.box.Quarter()
	nlevel := node.level + 1
	capacity := node.capacity

	node.children[TOP_LEFT] = NewNode(boxes[TOP_LEFT], nlevel, capacity)
	node.children[TOP_RIGHT] = NewNode(boxes[TOP_RIGHT], nlevel, capacity)
	node.children[BOTTOM_LEFT] = NewNode(boxes[BOTTOM_LEFT], nlevel, capacity)
	node.children[BOTTOM_RIGHT] = NewNode(boxes[BOTTOM_RIGHT], nlevel, capacity)
}

func (node *Node) GetPoints() []*shape.Point {
	return node.points
}

func NewQuadTree(x1, y1, x2, y2 float64, capacity int) (*QuadTree, error) {
	if x1 > x2 || y1 > y2 {
		return nil, &QuadTreeError{"Invalid Points for BoundaryBox construct"}
	}

	global := shape.NewBoundaryBox(
		shape.NewPoint(x1, y1),
		shape.NewPoint(x2, y2),
	)
	root := NewNode(global, 0, capacity)
	return &QuadTree{root, 0}, nil
}

func (node *Node) IsLeaf() bool {
	return node.children[TOP_RIGHT] == nil
}

func (node *Node) InsertPoint(point *shape.Point) bool {
	if !node.box.ContainsPoint(point) {
		return false
	}

	if node.length < node.capacity {
		node.points = append(node.points, point)
		node.length += 1
		return true
	}

	if node.IsLeaf() {
		node.SplitNode()
	}

	children := node.children
	return children[TOP_LEFT].InsertPoint(point) ||
		children[TOP_RIGHT].InsertPoint(point) ||
		children[BOTTOM_LEFT].InsertPoint(point) ||
		children[BOTTOM_RIGHT].InsertPoint(point)
}

func (qt *QuadTree) Insert(x, y float64) bool {
	p := shape.NewPoint(x, y)
	res := qt.root.InsertPoint(p)
	if res {
		qt.length += 1
	}

	return res
}

func (node *Node) GetPointsFromArea(area *shape.BoundaryBox) []*shape.Point {
	//we are not in valid node
	if node == nil || !node.box.Intersect(area) {
		return []*shape.Point{}
	}

	result := []*shape.Point{}
	log.Printf("AREA: %u\n", area)
	for _, point := range node.GetPoints() {
		if area.ContainsPoint(point) {
			log.Printf("contains: %u\n", point)
			result = append(result, point)
		}
	}

	if !node.IsLeaf() {
		children := node.children
		result = append(result, children[TOP_LEFT].GetPointsFromArea(area)...)
		result = append(result, children[TOP_RIGHT].GetPointsFromArea(area)...)
		result = append(result, children[BOTTOM_LEFT].GetPointsFromArea(area)...)
		result = append(result, children[BOTTOM_RIGHT].GetPointsFromArea(area)...)
	}

	return result
}

func (qt *QuadTree) GetPoints(x1, y1, x2, y2 float64) ([]*shape.Point, error) {
	if x1 > x2 || y1 > y2 {
		return nil, &QuadTreeError{"Invalid Points for BoundaryBox construct"}
	}

	area := shape.NewBoundaryBox(shape.NewPoint(x1, y1), shape.NewPoint(x2, y2))
	return qt.root.GetPointsFromArea(area), nil
}

func (qt *QuadTree) GetLength() int {
	return qt.length
}
