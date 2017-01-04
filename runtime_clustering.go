package main

import (
	"github.com/alldroll/spatial-index/clustering"
	"github.com/alldroll/spatial-index/geometry"
)

type RunTimeClustering struct {
	*shape.BoundaryBox
}

func NewRunTimeClustering(x1, y1, x2, y2 float64) *RunTimeClustering {
	return &RunTimeClustering{
		shape.NewBoundaryBox(shape.NewPoint(x1, y1), shape.NewPoint(x2, y2)),
	}
}

func (self *RunTimeClustering) GetClusters(points []*shape.Point, zoom int) []*shape.Cluster {
	bl, tr := self.GetBottomLeft(), self.GetTopRight()
	grid := cluster.NewGrid(bl.GetX(), bl.GetY(), tr.GetX(), tr.GetY(), zoom)
	grid.AddChunk(points)
	return grid.GetClusters()
}
