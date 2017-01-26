package main

import (
	"encoding/json"
	"github.com/alldroll/spatial-index"
	"github.com/alldroll/spatial-index/clustering"
	"github.com/alldroll/spatial-index/geometry"
	"log"
	"os"
	"strconv"
)

type Service struct {
	si *spatial_index.SpatialIndex
}

func NewService() *Service {
	si := spatial_index.NewSpatialIndex(23)
	service := &Service{si}
	service.loadMarkets("cmd/example/vl_points.json")
	return service
}

func (self *Service) RangeQuery(bounds *shape.BoundaryBox, quadKeys []string) []*shape.Cluster {
	var (
		allowed  []string
		clusters []*shape.Cluster
	)

	for _, quadKey := range quadKeys {
		cluster := self.si.GetCluster(quadKey)
		if cluster != nil && cluster.GetCount() < 500 {
			allowed = append(allowed, quadKey)
		} else if cluster != nil {
			clusters = append(clusters, cluster)
		}
	}

	points := self.si.RangeQuery(allowed)
	clusterBuilder := cluster.NewClusterBuilder(bounds, 0.1)
	for _, point := range points {
		clusterBuilder.AddPoint(point)
	}

	return append(clusterBuilder.GetClusters(), clusters...)
}

func (self *Service) loadMarkets(source string) {
	sourceT := struct {
		Points []struct{ Lat, Lng string }
	}{}

	file, _ := os.Open(source)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&sourceT)
	if err != nil {
		log.Fatal(err)
	}

	for _, point := range sourceT.Points {
		x, _ := strconv.ParseFloat(point.Lat, 64)
		y, _ := strconv.ParseFloat(point.Lng, 64)
		self.si.Insert(x, y)
	}
}
