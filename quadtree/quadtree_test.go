package quadtree

import (
    "testing"
    "log"
    "math/rand"
    "math"
)

func TestHudeData(t *testing.T) {
    qt, err := NewQuadTree(0, 0, 1, 1, 1. / math.Pow(4, 6))
    if (err != nil) {
        log.Fatal(err)
    }

    i := 0
    for i < 100000 {
        qt.Insert(rand.Float64(), rand.Float64())
        i += 1
    }

    points, err := qt.GetPoints(0.0, 0.0, 0.01, 0.01)
    if (err != nil) {
        log.Fatal(err)
    }

    t.Logf("COUNT IN RANGE %d", len(points))

    if (len(points) == 0) {
        t.Errorf("fail, point: %u", points)
    }
}

func TestSimple(t *testing.T) {
    qt, err := NewQuadTree(0, 0, 1, 1, 1. / math.Pow(4, 2))
    if (err != nil) {
        log.Fatal(err)
    }

    qt.Insert(0, 0.2)
    qt.Insert(0, 0.7)
    qt.Insert(0.9, 0.1)
    qt.Insert(0.5, 0.5)
    qt.Insert(0.1, 0.1)
    qt.Insert(0.3, 0.5)

    points, err := qt.GetPoints(0.48, 0.48, 0.49, 0.49)
    if (err != nil) {
        log.Fatal(err)
    }

    if (len(points) != 0) {
        t.Errorf("fail, point: %u", points)
    }

    points, err = qt.GetPoints(0, 0, 0.3, 0.3)

    if (err != nil) {
        log.Fatal(err)
    }

    if (len(points) == 0) {
        t.Errorf("fail, point: %u", points)
    }

    t.Logf("%u", points[0].EqualXY(0, 0.2))
}
