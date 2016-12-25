package cluster

import (
	"github.com/alldroll/spatial-index/geometry"
	"github.com/alldroll/spatial-index/quadtree"
)

type Cluster struct {
	center *shape.Point
	count  int
}

const (
	unclassified = 0
	noise        = 1
	busy         = 2
)

type Clustering struct {
	spartialIndex   *quadtree.QuadTree
	eps             float64
	minPoints       int
	classifications map[*shape.Point]byte
}

func (self *Cluster) GetCenter() *shape.Point {
	return self.center
}

func (self *Cluster) GetCount() int {
	return self.count
}

func NewClustering(quadtree *quadtree.QuadTree, db []*shape.Point, eps float64, minPoints int) *Clustering {
	for _, point := range db {
		quadtree.InsertPoint(point)
	}

	return &Clustering{
		quadtree,
		eps,
		minPoints,
		make(map[*shape.Point]byte),
	}
}

func (self *Clustering) DBScan(db []*shape.Point, eps float64, minPoints int) []*Cluster {
	var clusters []*Cluster

	for _, point := range db {
		if self.classifications[point] != unclassified {
			continue
		}

		neighborPoints := self.regionQuery(point)
		if len(neighborPoints) < minPoints {
			self.classifications[point] = noise
		} else {
			curCluster := &Cluster{point, 1}
			self.classifications[point] = busy
			self.expandCluster(curCluster, neighborPoints)
			clusters = append(clusters, curCluster)
		}
	}

	for point, status := range self.classifications {
		if status == noise {
			curCluster := &Cluster{point, 1}
			clusters = append(clusters, curCluster)
		}
	}

	return clusters
}

func (self *Clustering) expandCluster(cluster *Cluster, neighborPoints []*shape.Point) {
	for _, point := range neighborPoints {
		if self.classifications[point] == unclassified {
			neighborPointsForCur := self.regionQuery(point)
			if len(neighborPointsForCur) >= self.minPoints {
				neighborPoints = append(neighborPoints, neighborPointsForCur...)
			}
		}

		if self.classifications[point] != busy {
			cluster.count++
			self.classifications[point] = busy
		}
	}
}

func (self *Clustering) regionQuery(point *shape.Point) []*shape.Point {
	points, _ := self.spartialIndex.GetPoints(
		point.GetX()-self.eps,
		point.GetY()-self.eps,
		point.GetX()+self.eps,
		point.GetY()+self.eps,
	)

	/* if */
	return points
}
