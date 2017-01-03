package quadtree

import (
	"log"
	"math/rand"
	"testing"
)

func TestLength(t *testing.T) {
	qt, err := NewQuadTree(0, 0, 1, 1, 6)
	if err != nil {
		log.Fatal(err)
	}

	qt.Insert(0, 0.2)
	qt.Insert(0, 0.7)
	qt.Insert(0.9, 0.1)
	qt.Insert(0.5, 0.5)
	qt.Insert(0.1, 0.1)
	qt.Insert(0.3, 0.5)

	if qt.length != 6 {
		log.Fatal("Should be length 6")
	}

	qt.Insert(0, 0)

	if qt.length != 7 {
		log.Fatal("Should be length 6")
	}
}

func TestSimple(t *testing.T) {
	qt, err := NewQuadTree(0, 0, 1, 1, 2)
	if err != nil {
		log.Fatal(err)
	}

	qt.Insert(0, 0.2)
	qt.Insert(0, 0.7)
	qt.Insert(0.9, 0.1)
	qt.Insert(0.5, 0.5)
	qt.Insert(0.1, 0.1)
	qt.Insert(0.3, 0.5)
	qt.Insert(0, 0)

	points, err := qt.GetPoints(0, 0, 1, 1)

	if err != nil {
		log.Fatal(err)
	}

	if len(points) != 7 {
		t.Errorf("Should be 7 points, %u", points)
	}

	points, err = qt.GetPoints(0, 0, 0.2, 0.2)
	if err != nil {
		log.Fatal(err)
	}

	if len(points) != 3 {
		t.Errorf("Should be 3 points, %u", points)
	}
}

func TestShouldBeOnePoint(t *testing.T) {
	qt, err := NewQuadTree(0, 0, 1, 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	qt.Insert(0, 0)
	qt.Insert(0.1, 0.1)
	qt.Insert(0, 0.12)

	points, err := qt.GetPoints(0.025, 0.025, 0.125, 0.125)
	if err != nil {
		log.Fatal(err)
	}

	if len(points) != 1 {
		t.Errorf("Should be 1 point, %u", points)
	}
}

func BenchmarkGetPoints(b *testing.B) {
	qt, err := NewQuadTree(0, 0, 1, 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		a := 0 + rand.Float64()
		b := 0 + rand.Float64()
		qt.Insert(a, b)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		a := 0 + rand.Float64()
		b := 0 + rand.Float64()
		qt.GetPoints(a, a+0.2, b, b+0.2)
	}
}
