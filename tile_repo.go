package main

import (
	"encoding/json"
	"github.com/alldroll/spatial-index/geometry"
	"github.com/alldroll/spatial-index/tile_system"
	"github.com/alldroll/spatial-index/trie"
	"log"
	"os"
	"strconv"
	"sync/atomic"
)

type TileRepo struct {
	tr atomic.Value
}

type source2T struct {
	Points []struct{ Lat, Lng string }
}

func NewTileRepo(source string) *TileRepo {
	repo := &TileRepo{}
	repo.tr.Store(trie.NewQuadKeyTrie())
	repo.loadMarkets(source)
	return repo
}

func (self *TileRepo) RangeQueryQuadKeys(quadKeys []string) (string, []*shape.Point) {
	commonPrefix := LCP(quadKeys)

	var data []*shape.Point
	for _, quadKey := range quadKeys {
		d, _ := self.tr.Load().(*trie.QuadKeyTrie).RangeQuery([]byte(quadKey))
		data = append(data, d...)
	}

	return commonPrefix, data
}

func (self *TileRepo) GetCluster(quadKey string) *shape.Cluster {
	return self.tr.Load().(*trie.QuadKeyTrie).GetCluster([]byte(quadKey))
}

func (self *TileRepo) InsertPoint(lat, lng float64) {
	tileX, tileY := tile_system.LatLngToTileXY(lat, lng, 23)
	quadKey := tile_system.TileXYToQuadKey(tileX, tileY, 23)
	trie, err := self.tr.Load().(*trie.QuadKeyTrie).AddPoint(
		[]byte(quadKey),
		shape.NewPoint(lat, lng),
	)

	if err != nil {
		log.Fatal(err)
	}

	self.tr.Store(trie)
}

func (self *TileRepo) RemovePoint(x, y float64) {
	/* Implement me */
}

func (self *TileRepo) loadMarkets(source string) {
	sourceT := source2T{}

	file, _ := os.Open(source)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&sourceT)
	if err != nil {
		log.Fatal(err)
	}

	for _, point := range sourceT.Points {
		x, _ := strconv.ParseFloat(point.Lat, 64)
		y, _ := strconv.ParseFloat(point.Lng, 64)
		self.InsertPoint(x, y)
	}
}
