package btree

import (
	"bytes"
	"github.com/alldroll/quadtree/geometry"
)

type byteBuf bytes.Buffer

type DataStorage interface {
	Read(page int) byteBuf
	Write(page int, byteBuf *buf)
}

type BTree struct {
	root   *node
	height int
	degree int
	ds     *DataStorage
}

type Key interface {
	Less(then Key) bool
}

type node struct {
	leaf  bool
	keys  []Key
	btree *BTree
}

func NewBTree(ds *DataStorage, degree int) *BTree {
	if degree < 2 {
		panic("bad degree")
	}

	btree := &BTree{
		root:   nil,
		height: 0,
		degree: degree,
		ds:     ds,
	}

	node := btree.newNode()
	node.leaf = true
	btree.root = node
	/*btree.ds.Write()*/

	return btree
}

func (BTree *btree) Search(Key key) {

}

func (BTree *btree) newNode() *node {
	return &node{leaf: false, keys: []Key{}, btree: btree}
}
