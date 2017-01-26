package trie

import (
	"github.com/alldroll/spatial-index/geometry"
	"math/rand"
	"testing"
)

func TestRangeQuery(t *testing.T) {
	trie := NewQuadKeyTrie()

	quadKeys := [...][]byte{
		[]byte("1202300"),
		[]byte("1202201"),
		[]byte("1202202"),
		[]byte("1202312"),
		[]byte("0312313"),
		[]byte("0312330"),
		[]byte("1212202"),
		[]byte("1202301"),
		[]byte("120"),
		[]byte("12"),
	}

	for i, quadKey := range quadKeys {
		newTrie, err := trie.AddPoint(quadKey, shape.NewPoint(float64(i), float64(i)))
		if err != nil {
			t.Errorf(err.Error())
		}

		trie = newTrie
	}

	points, err := trie.RangeQuery([]byte("120"))
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(points) != 6 {
		t.Errorf(
			"Test Fail, expected {%d} len", 6,
		)
	}

	cluster := trie.GetCluster([]byte("120"))
	if 5 != cluster.GetCount() {
		t.Errorf(
			"Test Fail, expected {%d} len, got {%d}", 5, cluster.GetCount(),
		)
	}
}

func TestMissingRangeQuery(t *testing.T) {
	trie := NewQuadKeyTrie()

	quadKeys := [...][]byte{
		[]byte("1202300"),
		[]byte("1202201"),
		[]byte("1202202"),
		[]byte("1202312"),
		[]byte("0312313"),
		[]byte("0312330"),
		[]byte("1212202"),
		[]byte("1202301"),
		[]byte("120"),
		[]byte("12"),
	}

	for i, quadKey := range quadKeys {
		newTrie, err := trie.AddPoint(quadKey, shape.NewPoint(float64(i), float64(i)))
		if err != nil {
			t.Errorf(err.Error())
		}

		trie = newTrie
	}

	points, err := trie.RangeQuery([]byte("000"))
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(points) != 0 {
		t.Errorf(
			"Test Fail, expected {%d} len", 0,
		)
	}
}

func BenchmarkCreateTrie(b *testing.B) {
	var quadKeys [][]byte
	for i := 0; i < 10000; i++ {
		quadKeys = append(quadKeys, randBytes(23))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		trie := NewQuadKeyTrie()
		for i, quadKey := range quadKeys {
			trie, _ = trie.AddPoint(quadKey, shape.NewPoint(float64(i), float64(i)))
		}
	}
}

const letterBytes = "0123"

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return b
}
