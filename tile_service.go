package main

import (
	"github.com/alldroll/spatial-index/geometry"
	//"log"
	"github.com/alldroll/spatial-index/clustering"
	//"strconv"
)

type TileService struct {
	tr *TileRepo
}

func NewTileService(tr *TileRepo) *TileService {
	return &TileService{tr}
}

func (self *TileService) RangeQueryQuadKeys(bounds *shape.BoundaryBox, quadKeys []string) []*shape.Cluster {
	_, nodesData := self.tr.RangeQueryQuadKeys(quadKeys)
	clusterBuilder := cluster.NewClusterBuilder(bounds, 0.1)
	for _, nodeData := range nodesData {
		for _, point := range nodeData.GetData() {
			clusterBuilder.AddPoint(point.(*shape.Point))
		}
	}

	return clusterBuilder.GetClusters()
}
