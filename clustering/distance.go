package cluster

import (
	"github.com/alldroll/spatial-index/geometry"
	"math"
)

type ClusterBuilder struct {
	clusters    []*shape.Cluster
	bounds      *shape.BoundaryBox
	maxDistance float64
	factor      float64
}

func NewClusterBuilder(factor float64) *ClusterBuilder {
	return &ClusterBuilder{
		make([]*shape.Cluster, 0), nil, 0, factor,
	}
}

func (self *ClusterBuilder) AddPoint(point *shape.Point) {
	if self.bounds == nil {
		self.bounds = shape.NewBoundaryBox(point, point)
	} else if !self.bounds.ContainsPoint(point) {
		self.bounds = self.bounds.ExtendPoint(point)
		self.maxDistance = self.factor * boundsSize(self.bounds)
	}

	var nearest *shape.Cluster = nil
	for _, cluster := range self.clusters {
		if cluster.GetCenter().Equal(point) {
			nearest = cluster
			break
		}

		if self.maxDistance > 0 {
			distance := distance2Points(cluster.GetCenter(), point)
			if distance <= self.maxDistance && boundsSize(cluster.ExtendPoint(point)) <= self.maxDistance {
				nearest = cluster
			}
		}
	}

	if nearest != nil {
		nearest.AddPoint(point)
	} else {
		self.clusters = append(self.clusters, shape.NewCluster(point, 1))
	}
}

func (self *ClusterBuilder) GetClusters() []*shape.Cluster {
	return self.clusters
}

const (
	toRadians   = math.Pi / 180
	earthRadius = 6371.0087714
)

func boundsSize(bounds *shape.BoundaryBox) float64 {
	return distance2Points(
		bounds.GetBottomLeft(),
		bounds.GetTopRight(),
	)
}

func distance2Points(a, b *shape.Point) float64 {
	return haversinKilometers(
		a.GetX(),
		a.GetY(),
		b.GetX(),
		b.GetY(),
	)
}

func haversinKilometers(lat1, lng1, lat2, lng2 float64) float64 {
	x1, x2 := lat1*toRadians, lat2*toRadians
	h1 := 1 - math.Cos(x2-x1)
	h2 := 1 - math.Cos((lng2-lng1)*toRadians)
	h := h1 + math.Cos(x1)*math.Cos(x2)*h2

	hvr := math.Float64frombits(math.Float64bits(h) & 0xFFFFFFFFFFFFFFF8)

	return earthRadius * 2 * math.Asin(math.Min(1, math.Sqrt(hvr*0.5)))
}
