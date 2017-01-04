package main

import "github.com/alldroll/spatial-index/geometry"

type repo interface {
	RangeQuery(x1, y1, x2, y2 float64) []*shape.Point
	InsertPoint(x, y float64)
	RemovePoint(x, y float64)
}

type clustering interface {
	GetClusters(points []*shape.Point, zoom int) []*shape.Cluster
}

type Service struct {
	si repo
	cl clustering
}

func NewService(si repo, cl clustering) *Service {
	return &Service{si, cl}
}

func (self *Service) RangeQuery(x1, y1, x2, y2 float64, zoom int) []*shape.Cluster {
	points := self.si.RangeQuery(x1, y1, x2, y2)
	return self.cl.GetClusters(points, zoom)
}
