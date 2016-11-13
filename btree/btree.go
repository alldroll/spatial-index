package btree

import (
	"bytes"
)

type BTree struct {
	root   *node
	length uint
	degree uint
}

type Key interface {
	Less(than Key) bool
}

type node struct {
	numkeys  uint
	leaf     bool
	keys     []Key
	children []*node
	btree    *BTree
}

func NewBTree(degree int) *BTree {
	if degree < 2 {
		panic("bad degree")
	}

	return &BTree{
		degree: degree,
	}
}

func (btree *BTree) Search(Key key) Key {
	if btree.root != nil {
		node, index := btree.root.find(key)
		if node {
			return key
		}
	}

	return nil
}

func (btree *BTree) Insert(Key key) Key {
	if btree.root == nil {
		return nil
	}

	return btree.root.insert(key)
}

func (n *node) insert(Key key) *node {
	//if len(n.keys) <
}

func (btree *BTree) splitChild(n *node, uint index) {
	t = btree.degree
	y = n.children[index]
	z := btree.newNode()
	z.leaf = y.leaf

	z.keys = append(z.keys, y.keys[index+1:]...)
	/*
		for j := 1; j < t-1; j++ {
			z.keys[j] = y.keys[j+t]
		}
	*/

	if !y.leaf {
		z.children = append(z.children, y.children[index+1:]...)
		/*
			for j := 1; j < t; j++ {
				z.children[j] = y.children[j+t]
			}
		*/
	}

	x.keys = append(x.keys[:index], append(y.keys[t-1], x.keys[:index]...))
	x.children = append(x.children[:index+1], append(z, x.children[:index+1]...))
}

func (node *n) search(Key key) (*node, uint) {
	for index, skey := range n.keys {
		if !skey.Less(key) {
			break
		}
	}

	if index < len(n.keys) && !key.Less(skey) {
		return n.children[index], index
	}

	if n.leaf {
		return nil, nil
	}

	return n.children[index].search(key)
}

func (BTree *btree) newNode() *node {
	t = btree.degree
	return &node{leaf: false, keys: []Key{}, children: []*node{}, btree: btree}
}
