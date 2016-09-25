package cluster

import (
	"github.com/alldroll/quadtree/geometry"
	"github.com/alldroll/quadtree/quadtree"
)

/* declare */
func DistanceBetweenPoints(a *shape.Point, b *shape.Point) float32 /*?*/

type Clustering struct {
	qt *quadtree.QuadTree
}

type Cluster struct {
	center *shape.Point
	count  uint
}

func (Clustering *self) DBScan(db []*shape.Point, eps float64, minPoints uint) ([]*shape.Point, []*Cluster) {
	var (
		spoints  []*shape.Point
		clusters []*Cluster
		visited  = make([]bool, len(db))
	)

	for i, point := range db {
		if visited[i] {
			continue
		}

		visited[i] = true
		neighborPoints = qt.regionQuery(point, eps)
		if len(neighborPoints) > minPoints {
			points = append(points, point...)
		} else {
			/* */
		}
	}
}
