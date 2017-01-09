package trie

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

func TestTrieLookup(t *testing.T) {
	trie := NewTrie()

	words := [...]string{
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

	expected := []string{
		"1202201",
		"1202202",
		"1202300",
		"120",
		"1202301",
		"1202312",
	}

	for i, word := range words {
		err := trie.AddWord(word, i)
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	nodesData, err := trie.Lookup("120")
	if err != nil {
		t.Errorf(err.Error())
	}
	var actual []string

	for _, nodeData := range nodesData {
		actual = append(actual, nodeData.GetWord())
	}

	sort.Sort(sort.StringSlice(expected))
	sort.Sort(sort.StringSlice(actual))

	if fmt.Sprintf("%v", expected) != fmt.Sprintf("%v", actual) {
		t.Errorf(
			"Test Fail, expected {%v}, got {%v}",
			expected,
			actual,
		)
	}
}

func BenchmarkGetPrefixes(b *testing.B) {
	trie := NewTrie()

	for i := 0; i < 200000; i++ {
		quadKey := randStringBytes(23)
		trie.AddWord(quadKey, i)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		prefixLen := rand.Intn(23)
		prefix := randStringBytes(prefixLen)
		b.StartTimer()
		trie.Lookup(prefix)
	}
}

const letterBytes = "0123"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
