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

func (self *TileService) RangeQueryTiles(x1, y1, x2, y2 int, zoom int) []*shape.Cluster {
	_, nodesData := self.tr.RangeQueryTiles(x1, y1, x2, y2, uint(zoom))
	clusterBuilder := cluster.NewClusterBuilder(0.1)
	for _, nodeData := range nodesData {
		for _, point := range nodeData.GetData() {
			clusterBuilder.AddPoint(point.(*shape.Point))
		}
	}

	return clusterBuilder.GetClusters()
}
