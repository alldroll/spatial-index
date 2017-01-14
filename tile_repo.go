package main

import (
	"encoding/json"
	"github.com/alldroll/spatial-index/geometry"
	"github.com/alldroll/spatial-index/tile_system"
	"github.com/alldroll/spatial-index/trie"
	"log"
	"os"
	"strconv"
)

type TileRepo struct {
	tr *trie.Trie
}

type source2T struct {
	Points []struct {
		Lat string
		Lng string
	}
}

func NewTileRepo(source string) *TileRepo {
	repo := &TileRepo{trie.NewTrie()}
	repo.loadMarkets(source)
	return repo
}

/*
func (self *TileRepo) RangeQuery(x1, y1, x2, y2 float64, zoom uint) (string, []*trie.NodeData) {
	tileMinX, tileMinY := tile_system.LatLngToTileXY(x1, y1, zoom)
	tileMaxX, tileMaxY := tile_system.LatLngToTileXY(x2, y2, zoom)
	return self.RangeQueryTiles(tileMinX, tileMinY, tileMaxX, tileMaxY, zoom)
}

func (self *TileRepo) RangeQueryTiles(tileMinX, tileMinY, tileMaxX, tileMaxY int, zoom uint) (string, []*trie.NodeData) {
	var quadKeys []string
	for i, len1 := 0, (tileMaxX - tileMinX); i <= len1; i++ {
		for j, len2 := 0, (tileMinY - tileMaxY); j <= len2; j++ {
			quadKey := tile_system.TileXYToQuadKey(tileMinX+i, tileMaxY+j, zoom)
			quadKeys = append(quadKeys, quadKey)
		}
	}

	return self.RangeQueryQuadKeys(quadKeys)
}
*/

func (self *TileRepo) RangeQueryQuadKeys(quadKeys []string) (string, []*trie.NodeData) {
	commonPrefix := LCP(quadKeys)

	var data []*trie.NodeData
	for _, quadKey := range quadKeys {
		d, _ := self.tr.Lookup(quadKey)
		data = append(data, d...)
	}

	return commonPrefix, data
}

func (self *TileRepo) InsertPoint(lat, lng float64) {
	tileX, tileY := tile_system.LatLngToTileXY(lat, lng, 23)
	quadKey := tile_system.TileXYToQuadKey(tileX, tileY, 23)
	self.tr.AddWord(quadKey, shape.NewPoint(lat, lng))
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
