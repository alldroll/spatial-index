package main

import (
	"github.com/alldroll/spatial-index/geometry"
	//"log"
	"strconv"
)

type TileService struct {
	tr *TileRepo
}

func NewTileService(tr *TileRepo) *TileService {
	return &TileService{tr}
}

func (self *TileService) RangeQuery(x1, y1, x2, y2 float64, zoom int) []*shape.Cluster {
	commonPrefix, nodesData := self.tr.RangeQuery(x1, y1, x2, y2, uint(zoom))
	len := len(commonPrefix)
	var grid [64]*shape.Cluster

	//log.Printf("COMMON %s\n", commonPrefix)

	for _, nodeData := range nodesData {
		quadKey := nodeData.GetWord()

		a, b, c := string(quadKey[len+3]), string(quadKey[len+4]), string(quadKey[len+5])
		ai, _ := strconv.Atoi(a)
		bi, _ := strconv.Atoi(b)
		ci, _ := strconv.Atoi(c)

		i := ai + (bi * 4) + (ci * 16)

		//log.Printf("%s, %d, %d, %d, z %d\n", quadKey, a, b, i, zoom)
		if grid[i] == nil {
			grid[i] = shape.NewCluster(shape.NewPoint(0, 0), 0)
		}

		for _, point := range nodeData.GetData() {
			grid[i].AddPoint(point.(*shape.Point))
		}
	}

	box := shape.NewBoundaryBox(shape.NewPoint(x1, y1), shape.NewPoint(x2, y2))
	var result []*shape.Cluster
	for _, cell := range grid {
		if cell != nil && cell.GetCount() > 0 && box.ContainsPoint(cell.GetCenter()) {
			result = append(result, cell)
		}
	}

	return result
}
