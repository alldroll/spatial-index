package trie

import (
	"errors"
	"github.com/alldroll/spatial-index/geometry"
)

type QuadKeyTrie struct {
	root *tile
}

type tile struct {
	quadKey []byte
	data    []*shape.Point
	edges   [4]*tile
	cluster *shape.Cluster
}

func NewQuadKeyTrie() *QuadKeyTrie {
	return &QuadKeyTrie{newEmptyTile()}
}

func (self *QuadKeyTrie) AddPoint(quadKey []byte, point *shape.Point) (*QuadKeyTrie, error) {
	if len(quadKey) == 0 {
		return nil, errors.New("empty quadKey")
	}

	newTrie := NewQuadKeyTrie()
	current := self.root.Copy()

	for iter := 0; iter < len(quadKey); iter++ {
		if iter == 0 {
			newTrie.root = current
		}

		if current.cluster == nil {
			current.cluster = shape.NewCluster(point, 1)
		} else {
			current.cluster.AddPoint(point)
		}

		k := quadKey[iter] - 48
		if current.edges[k] == nil {
			current.edges[k] = newEmptyTile()
		} else {
			current.edges[k] = current.edges[k].Copy()
		}

		current = current.edges[k]

		if iter == len(quadKey)-1 {
			if current.quadKey == nil {
				current.quadKey = quadKey
			}

			current.data = append(current.data, point)
		}
	}

	return newTrie, nil
}

func (self *QuadKeyTrie) RangeQuery(prefix []byte) ([]*shape.Point, error) {
	var data []*shape.Point
	if len(prefix) == 0 {
		return data, errors.New("empty prefix")
	}

	self.root.lookup(prefix, &data)
	if len(data) == 0 {
		return data, nil
	}

	copy := make([]*shape.Point, len(data))
	for i, point := range data {
		copy[i] = point.Copy()
	}

	return copy, nil
}

func (self *QuadKeyTrie) GetCluster(quadKey []byte) *shape.Cluster {
	current := self.root
	for i := 0; i < len(quadKey) && current != nil; i++ {
		k := quadKey[i] - 48
		current = current.edges[k]
	}

	if current == nil {
		return nil
	}

	return current.cluster
}

func newEmptyTile() *tile {
	return &tile{
		nil, make([]*shape.Point, 0), [4]*tile{nil, nil, nil, nil}, nil,
	}
}

func (self *tile) Copy() *tile {
	newTile := newEmptyTile()

	if self.cluster != nil {
		newTile.cluster = self.cluster.Copy()
	}

	for i, edge := range self.edges {
		newTile.edges[i] = edge
	}

	newTile.data = make([]*shape.Point, len(self.data))
	copy(newTile.data, self.data)
	if self.quadKey != nil {
		newTile.quadKey = make([]byte, len(self.quadKey))
		copy(newTile.quadKey, self.quadKey)
	}

	return newTile
}

func (self *tile) lookup(prefix []byte, data *[]*shape.Point) {
	if self.quadKey != nil && len(prefix) == 0 {
		*data = append(*data, self.data...)
	}

	if len(prefix) > 0 {
		k := prefix[0] - 48
		if self.edges[k] != nil {
			self.edges[k].lookup(prefix[1:], data)
		}
	} else {
		for _, edge := range self.edges {
			if edge != nil {
				edge.lookup(prefix, data)
			}
		}
	}
}
