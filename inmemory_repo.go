package main

import (
	"encoding/json"
	"github.com/alldroll/spatial-index/geometry"
	"github.com/alldroll/spatial-index/quadtree"
	"log"
	"os"
	"strconv"
)

type InMemoryRepo struct {
	qt *quadtree.QuadTree
}

type SourceT struct {
	Points []struct {
		Lat string
		Lng string
	}
}

func NewInMemoryRepo(x1, y1, x2, y2 float64, source string) *InMemoryRepo {
	qt, err := quadtree.NewQuadTree(x1, y1, x2, y2, 20)
	if err != nil {
		log.Fatal(err)
	}

	repo := &InMemoryRepo{qt}
	repo.loadMarkets(source)
	return repo
}

func (self *InMemoryRepo) RangeQuery(x1, y1, x2, y2 float64) []*shape.Point {
	points, err := self.qt.GetPoints(x1, y1, x2, y2)
	if err != nil {
		log.Fatal(err)
	}

	return points
}

func (self *InMemoryRepo) InsertPoint(x, y float64) {
	self.qt.Insert(x, y)
}

func (self *InMemoryRepo) RemovePoint(x, y float64) {
	/* Implement me */
}

func (self *InMemoryRepo) loadMarkets(source string) {
	sourceT := SourceT{}
	file, _ := os.Open(source)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&sourceT)
	if err != nil {
		log.Fatal(err)
	}

	for _, point := range sourceT.Points {
		x, _ := strconv.ParseFloat(point.Lat, 64)
		y, _ := strconv.ParseFloat(point.Lng, 64)
		self.qt.Insert(x, y)
	}
}
