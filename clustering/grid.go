package cluster

import (
	"github.com/alldroll/spatial-index/geometry"
)

type Grid struct {
	root     *node
	height   int
	clusters []*Cluster
}

type node struct {
	box      *shape.BoundaryBox
	children [totalChild]*node
	level    int
	cluster  *Cluster
}

type Cluster struct {
	length int
	center *shape.Point
}

const (
	topLeft     int = 0
	topRight        = 1
	bottomLeft      = 2
	bottomRight     = 3
	totalChild      = 4
)

func NewGrid(x1, y1, x2, y2 float64, height int) *Grid {
	if x1 > x2 || y1 > y2 {
		return nil
	}

	global := shape.NewBoundaryBox(
		shape.NewPoint(x1, y1),
		shape.NewPoint(x2, y2),
	)
	root := newNode(global, 0)
	grid := &Grid{root, height, make([]*Cluster, 0)}

	grid.splitNodeUntil(root)
	return grid
}

func (self *Grid) AddChunk(chunk []*shape.Point) {
	for _, point := range chunk {
		self.root.insertPoint(point)
	}
}

func (self *Grid) GetClusters() []*Cluster {
	res := make([]*Cluster, 0)
	for _, cluster := range self.clusters {
		if cluster.length > 0 {
			res = append(res, cluster)
		}
	}

	return res
}

func (self *Cluster) GetCenter() *shape.Point {
	return shape.NewPoint(self.center.GetX()/float64(self.length), self.center.GetY()/float64(self.length))
}

func (self *Cluster) GetLength() int {
	return self.length
}

func newNode(box *shape.BoundaryBox, level int) *node {
	return &node{
		box,
		[totalChild]*node{nil, nil, nil, nil},
		level,
		nil,
	}
}

func (self *node) insertPoint(point *shape.Point) bool {
	if !self.box.ContainsPoint(point) {
		return false
	}

	if self.cluster != nil {
		self.cluster.length++
		self.cluster.center.SetX(point.GetX() + self.cluster.center.GetX())
		self.cluster.center.SetY(point.GetY() + self.cluster.center.GetY())
		return true
	}

	children := self.children
	return children[topLeft].insertPoint(point) ||
		children[topRight].insertPoint(point) ||
		children[bottomLeft].insertPoint(point) ||
		children[bottomRight].insertPoint(point)
}

func (self *Grid) splitNodeUntil(curNode *node) {
	if self.height <= curNode.level {
		cluster := &Cluster{0, shape.NewPoint(0, 0)}
		self.clusters = append(self.clusters, cluster)
		curNode.cluster = cluster
		return
	}

	curNode.splitNode()
	for i := 0; i < totalChild; i++ {
		self.splitNodeUntil(curNode.children[i])
	}
}

func (self *node) splitNode() {
	boxes := self.box.Quarter()
	nlevel := self.level + 1
	for i := 0; i < totalChild; i++ {
		self.children[i] = newNode(boxes[i], nlevel)
	}
}
