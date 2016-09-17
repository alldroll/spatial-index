package quadtree

import (
	"log"
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

	if len(points) != 4 {
		t.Errorf("Should be 4 points, %u", points)
	}
}
