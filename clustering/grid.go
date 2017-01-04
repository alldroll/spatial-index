package cluster

import (
	"github.com/alldroll/spatial-index/geometry"
)

const (
	totalChild = 4
)

type Grid struct {
	root     *cell
	height   int
	clusters []*shape.Cluster
}

type cell struct {
	*shape.BoundaryBox
	children [totalChild]*cell
	level    int
	cluster  *shape.Cluster
}

func NewGrid(x1, y1, x2, y2 float64, height int) *Grid {
	if x1 > x2 || y1 > y2 {
		return nil
	}

	global := shape.NewBoundaryBox(
		shape.NewPoint(x1, y1),
		shape.NewPoint(x2, y2),
	)
	root := newCell(global, 0)
	grid := &Grid{root, height, make([]*shape.Cluster, 0)}

	grid.splitCellUntil(root)
	return grid
}

func (self *Grid) AddChunk(chunk []*shape.Point) {
	for _, point := range chunk {
		self.root.insertPoint(point)
	}
}

func (self *Grid) GetClusters() []*shape.Cluster {
	res := make([]*shape.Cluster, 0)
	for _, cluster := range self.clusters {
		if cluster.GetCount() > 0 {
			res = append(res, cluster)
		}
	}

	return res
}

func newCell(box *shape.BoundaryBox, level int) *cell {
	return &cell{
		box,
		[totalChild]*cell{nil, nil, nil, nil},
		level,
		nil,
	}
}

func (self *cell) insertPoint(point *shape.Point) bool {
	if !self.ContainsPoint(point) {
		return false
	}

	if self.cluster != nil {
		self.cluster.AddPoint(point)
		return true
	}

	success := false
	for i := 0; i < totalChild && !success; i++ {
		success = self.children[i].insertPoint(point)
	}

	return success
}

func (self *Grid) splitCellUntil(curCell *cell) {
	if self.height <= curCell.level {
		cluster := shape.NewCluster(shape.NewPoint(0, 0), 0)
		self.clusters = append(self.clusters, cluster)
		curCell.cluster = cluster
		return
	}

	curCell.splitCell()
	for i := 0; i < totalChild; i++ {
		self.splitCellUntil(curCell.children[i])
	}
}

func (self *cell) splitCell() {
	boxes := self.Quarter()
	nlevel := self.level + 1
	for i := 0; i < totalChild; i++ {
		self.children[i] = newCell(boxes[i], nlevel)
	}
}
