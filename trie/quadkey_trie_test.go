package trie

import (
	"github.com/alldroll/spatial-index/geometry"
	"testing"
)

func TestRangeQuery(t *testing.T) {
	trie := NewQuadKeyTrie()

	quadKeys := [...]string{
		"1202300",
		"1202201",
		"1202202",
		"1202312",
		"0312313",
		"0312330",
		"1212202",
		"1202301",
		"120",
		"12",
	}

	for i, quadKey := range quadKeys {
		err := trie.AddPoint(quadKey, shape.NewPoint(float64(i), float64(i)))
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	points, err := trie.RangeQuery("120")
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(points) != 6 {
		t.Errorf(
			"Test Fail, expected {%d} len", 6,
		)
	}

	cluster := trie.GetCluster("120")
	if 5 != cluster.GetCount() {
		t.Errorf(
			"Test Fail, expected {%d} len, got {%d}", 5, cluster.GetCount(),
		)
	}
}
