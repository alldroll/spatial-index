package cluster

import (
	"github.com/alldroll/spatial-index/geometry"
	"math/rand"
	"testing"
)

func BenchmarkClustering(b *testing.B) {
	k := 200000

	b.StopTimer()
	points := make([]*shape.Point, k)
	for i := 0; i < k; i++ {
		a := 0 + rand.Float64()
		b := 0 + rand.Float64()
		points[i] = shape.NewPoint(a, b)
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		grid := NewGrid(0, 0, 1, 1, 3)
		grid.AddChunk(points)
		grid.GetClusters()
	}
}
