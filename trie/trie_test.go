package trie

import (
	"fmt"
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

	for _, word := range words {
		err := trie.AddWord(word)
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	actual, err := trie.Lookup("120")
	if err != nil {
		t.Errorf(err.Error())
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
