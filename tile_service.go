package main

import (
	"github.com/alldroll/spatial-index/clustering"
	"github.com/alldroll/spatial-index/geometry"
)

type TileService struct {
	tr *TileRepo
}

func NewTileService(tr *TileRepo) *TileService {
	return &TileService{tr}
}

func (self *TileService) RangeQueryQuadKeys(bounds *shape.BoundaryBox, quadKeys []string) []*shape.Cluster {
	var (
		allowed  []string
		clusters []*shape.Cluster
	)

	for _, quadKey := range quadKeys {
		cluster := self.tr.GetCluster(quadKey)
		if cluster != nil && cluster.GetCount() < 500 {
			allowed = append(allowed, quadKey)
		} else if cluster != nil {
			clusters = append(clusters, cluster)
		}
	}

	_, points := self.tr.RangeQueryQuadKeys(allowed)
	clusterBuilder := cluster.NewClusterBuilder(bounds, 0.1)
	for _, point := range points {
		clusterBuilder.AddPoint(point)
	}

	return append(clusterBuilder.GetClusters(), clusters...)
}
