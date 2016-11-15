package btree

import (
	"testing"
)

type myKey int

func (m myKey) Less(than Key) bool {
	return m < than.(myKey)
}

func TestImplementMe(t *testing.T) {
	btree := NewBTree(2)
	if btree == nil {
		t.Errorf("TestFail, btree is nil")
	}

	m := [10]myKey{1, 2, 5, 4, 3, 6, 10, 7, -1, 20}
	for _, v := range m {
		btree.Insert(v)
	}

	for _, v := range m {
		if btree.Search(v) != v {
			t.Errorf("TestFail, expected: %d", v)
		}
	}

	btree.Draw()
}
