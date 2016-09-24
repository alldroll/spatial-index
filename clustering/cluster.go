package cluster

import (
	"github.com/alldroll/quadtree/geometry"
)

/* declare */
func DistanceBetweenPoints(a *shape.Point, b *shape.Point) float32 /*?*/

type Cluster struct {
	center *shape.Point
	points []*shape.Point
	bounds *shape.BoundaryBox
}

type PointClusterer struct {
	clusters *[]Clusters
}

/*
 */
func (Cluster *self) addPoint(point *shape.Point) bool {

}

func Optics() {}
