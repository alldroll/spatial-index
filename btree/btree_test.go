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

	var m myKey = 1
	btree.Insert(m)

	actual := btree.Search(m)
	if actual != m {
		t.Errorf("TestFail, expected: %d", m)
	}
}
