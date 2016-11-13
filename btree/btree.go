package btree

type BTree struct {
	root   *node
	length int
	degree int
}

type Key interface {
	Less(than Key) bool
}

type node struct {
	numkeys  int
	leaf     bool
	keys     []Key
	children []*node
	/*btree    *BTree*/
}

func NewBTree(degree int) *BTree {
	if degree < 2 {
		panic("bad degree")
	}

	return &BTree{
		degree: degree,
	}
}

func (btree *BTree) Search(key Key) Key {
	if btree.root != nil {
		node, _ := btree.root.search(key)
		if node != nil {
			return key
		}
	}

	return nil
}

func (btree *BTree) Insert(key Key) Key {
	maxItems := btree.maxItems()

	if btree.root == nil {
		btree.root = btree.newNode()
		btree.root.leaf = true
	}

	r := btree.root
	if len(r.keys) < maxItems {
		return btree.insertNotFull(r, key)
	}

	s := btree.newNode()
	btree.root = s
	s.children = []*node{r}
	btree.splitChild(s, 0)

	return btree.insertNotFull(s, key)
}

func (btree *BTree) insertNotFull(n *node, key Key) Key {
	maxItems := btree.maxItems()
	i := len(n.keys) - 1
	if n.leaf {
		for i >= 0 && key.Less(n.keys[i]) {
			n.keys[i+1] = n.keys[i]
			i--
		}

		n.keys = append(n.keys, key)
		return key
	}

	for ; i >= 0 && key.Less(n.keys[i]); i-- {
	}

	i++
	if len(n.children[i].children) == maxItems {
		btree.splitChild(n, i)
		if n.keys[i].Less(key) {
			i++
		}
	}

	return btree.insertNotFull(n.children[i], key)
}

func (btree *BTree) maxItems() int {
	return btree.degree*2 - 1
}

func (n *node) search(key Key) (*node, int) {
	var (
		index int = 0
		skey  Key = nil
	)

	for i, sk := range n.keys {
		index, skey = i, sk
		if !skey.Less(key) {
			break
		}
	}

	if index < len(n.keys) && !key.Less(skey) {
		return n, index
	}

	if n.leaf {
		return nil, -1
	}

	return n.children[index].search(key)
}

func (btree *BTree) splitChild(n *node, index int) {
	t := btree.degree
	y := n.children[index]
	z := btree.newNode()
	z.leaf = y.leaf

	z.keys = append(z.keys, y.keys[index+1:]...)
	if !y.leaf {
		z.children = append(z.children, y.children[index+1:]...)
	}

	n.keys = append(n.keys[:index], append([]Key{y.keys[t-1]}, n.keys[:index]...)...)
	n.children = append(n.children[:index+1], append([]*node{z}, n.children[:index+1]...)...)
}

func (btree *BTree) newNode() *node {
	/*t := btree.degree*/
	return &node{leaf: false, keys: []Key{}, children: []*node{}}
}
