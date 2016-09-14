package quadtree

import (
    "github.com/alldroll/quadtree/geometry"
)

/**/
type Part int
const (
    TOP_LEFT Part = 0
    TOP_RIGHT = 1
    BOTTOM_LEFT = 2
    BOTTOM_RIGHT = 3
    NODE_CAPACITY = 40
)

/**/
type Node struct {
    box *shape.BoundaryBox
    children [4]*Node
    level int
}

/**/
type QuadTree struct {
    root *Node
    length int
}

/**/
type QuadTreeError struct {
    msg string
}

func (e *QuadTreeError) Error() string {
    return e.msg
}

func NewNode(box *shape.BoundaryBox, level int) *Node {
    return &Node{ box, [4]*Node{nil, nil, nil, nil}, level }
}

func DivideNodeUntil(node *Node, minArea float64) *Node {
    if (node == nil || node.box.Area() < minArea) {
        return nil
    }

    boxes := node.box.Quarter()
    nlevel := node.level + 1

    node.children[TOP_LEFT] = DivideNodeUntil(
        NewNode(boxes[TOP_LEFT], nlevel),
        minArea,
    )

    node.children[TOP_RIGHT] = DivideNodeUntil(
        NewNode(boxes[TOP_RIGHT], nlevel),
        minArea,
    )

    node.children[BOTTOM_LEFT] = DivideNodeUntil(
        NewNode(boxes[BOTTOM_LEFT], nlevel),
        minArea,
    )

    node.children[BOTTOM_RIGHT] = DivideNodeUntil(
        NewNode(boxes[BOTTOM_RIGHT], nlevel),
        minArea,
    )

    return node
}

func NewQuadTree(x1, y1, x2, y2 float64, minArea float64) (*QuadTree, error) {
    if (x1 > x2 || y1 > y2) {
        return nil, &QuadTreeError{"Invalid Points for BoundaryBox construct"}
    }
    global := shape.NewBoundaryBox(shape.NewPoint(x1, y1), shape.NewPoint(x2, y2))
    root := DivideNodeUntil(NewNode(global, 0), minArea)
    return &QuadTree{root, 0}, nil
}

func InsertPoint(cur* Node, point *shape.Point) bool {
    if (cur == nil || !cur.box.ContainsPoint(point)) {
        return false
    }

    children := cur.children
    if (children == [4]*Node{nil, nil, nil, nil}) {
        cur.box.AppendPoint(point)
        return true
    }

    if (cur.box.GetPointsCount() < NODE_CAPACITY) {
        cur.box.AppendPoint(point)
    }

    return InsertPoint(children[TOP_LEFT], point) ||
        InsertPoint(children[TOP_RIGHT], point) ||
        InsertPoint(children[BOTTOM_LEFT], point) ||
        InsertPoint(children[BOTTOM_RIGHT], point)
}

func (qt *QuadTree) Insert(x, y float64) bool {
    p := shape.NewPoint(x, y)
    return InsertPoint(qt.root, p)
}

func GetPointsFromArea(cur *Node, area *shape.BoundaryBox) []*shape.Point {
    //we are not in valid node
    if (cur == nil || !cur.box.Intersect(area)) {
        return []*shape.Point{}
    }

    if (!cur.box.ContainsBox(area)) {
        return cur.box.GetPoints()
    }

    //if this is leaf return points
    if (cur.children == [4]*Node{nil, nil, nil, nil}) {
        return cur.box.GetPoints()
    }

    result := GetPointsFromArea(cur.children[TOP_LEFT], area)
    result = append(result, GetPointsFromArea(cur.children[TOP_RIGHT], area)...)
    result = append(result, GetPointsFromArea(cur.children[BOTTOM_LEFT], area)...)
    result = append(result, GetPointsFromArea(cur.children[BOTTOM_RIGHT], area)...)

    return result
}

func (qt *QuadTree) GetPoints(x1, y1, x2, y2 float64) ([]*shape.Point, error) {
    if (x1 > x2 || y1 > y2) {
        return nil, &QuadTreeError{"Invalid Points for BoundaryBox construct"}
    }

    area := shape.NewBoundaryBox(shape.NewPoint(x1, y1), shape.NewPoint(x2, y2))
    return GetPointsFromArea(qt.root, area), nil
}
