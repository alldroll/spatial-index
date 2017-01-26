package spatial_index

import (
	"github.com/alldroll/spatial-index/geometry"
	"github.com/alldroll/spatial-index/tile_system"
	"github.com/alldroll/spatial-index/trie"
	"log"
	"sync/atomic"
)

type SpatialIndex struct {
	tr      atomic.Value
	maxZoom uint
}

func NewSpatialIndex(maxZoom uint) *SpatialIndex {
	spatialIndex := &SpatialIndex{}
	spatialIndex.tr.Store(trie.NewQuadKeyTrie())
	spatialIndex.maxZoom = maxZoom
	return spatialIndex
}

func (self *SpatialIndex) RangeQuery(quadKeys []string) []*shape.Point {
	var data []*shape.Point
	for _, quadKey := range quadKeys {
		d, _ := self.tr.Load().(*trie.QuadKeyTrie).RangeQuery([]byte(quadKey))
		data = append(data, d...)
	}

	return data
}

func (self *SpatialIndex) GetCluster(quadKey string) *shape.Cluster {
	return self.tr.Load().(*trie.QuadKeyTrie).GetCluster([]byte(quadKey))
}

func (self *SpatialIndex) Insert(lat, lng float64) {
	tileX, tileY := tile_system.LatLngToTileXY(lat, lng, self.maxZoom)
	quadKey := tile_system.TileXYToQuadKey(tileX, tileY, self.maxZoom)
	trie, err := self.tr.Load().(*trie.QuadKeyTrie).AddPoint(
		[]byte(quadKey),
		shape.NewPoint(lat, lng),
	)

	if err != nil {
		log.Fatal(err)
	}

	self.tr.Store(trie)
}

func (self *SpatialIndex) RemovePoint(x, y float64) {
	panic("Not implemented yet")
}
