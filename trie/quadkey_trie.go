package trie

import (
	"errors"
	"github.com/alldroll/spatial-index/geometry"
	"strconv"
)

type QuadKeyTrie struct {
	root   *tile
	length int
}

type tile struct {
	quadKey string
	data    []*shape.Point
	edges   [4]*tile
	cluster *shape.Cluster
}

func NewQuadKeyTrie() *QuadKeyTrie {
	return &QuadKeyTrie{
		newTile(),
		0,
	}
}

func (self *QuadKeyTrie) AddPoint(quadKey string, point *shape.Point) error {
	if len(quadKey) == 0 {
		return errors.New("empty quadKey")
	}

	self.root.addPoint(quadKey, 0, point)
	return nil
}

func (self *QuadKeyTrie) RangeQuery(prefix string) ([]*shape.Point, error) {
	var data []*shape.Point
	if len(prefix) == 0 {
		return data, errors.New("empty prefix")
	}

	self.root.lookup(prefix, 0, &data)
	if len(data) == 0 {
		return data, nil
	}

	copy := make([]*shape.Point, len(data))
	for i, point := range data {
		copy[i] = point.Copy()
	}

	return copy, nil
}

func (self *QuadKeyTrie) GetCluster(quadKey string) *shape.Cluster {
	cluster := self.root.getCluster(quadKey, 0)
	if cluster != nil {
		cluster = cluster.Copy()
	}

	return cluster
}

func newTile() *tile {
	return &tile{
		"", make([]*shape.Point, 0), [4]*tile{nil, nil, nil, nil}, nil,
	}
}

func (self *tile) addPoint(quadKey string, iter int, point *shape.Point) {
	if len(quadKey) <= iter {
		self.quadKey = quadKey
		self.data = append(self.data, point)
		return
	}

	if self.cluster == nil {
		self.cluster = shape.NewCluster(point, 1)
	} else {
		self.cluster.AddPoint(point)
	}

	c := string(quadKey[iter])
	k, _ := strconv.Atoi(c)
	if self.edges[k] == nil {
		self.edges[k] = newTile()
	}

	iter++
	self.edges[k].addPoint(quadKey, iter, point)
}

func (self *tile) lookup(prefix string, iter int, data *[]*shape.Point) {
	if len(self.quadKey) > 0 && len(prefix) == iter {
		*data = append(*data, self.data...)
	}

	if len(prefix) > iter {
		c := string(prefix[iter])
		k, _ := strconv.Atoi(c)
		iter++
		if self.edges[k] != nil {
			self.edges[k].lookup(prefix, iter, data)
		}
	} else {
		for _, edge := range self.edges {
			if edge != nil {
				edge.lookup(prefix, iter, data)
			}
		}
	}
}

func (self *tile) getCluster(quadKey string, iter int) *shape.Cluster {
	if len(quadKey) == iter {
		return self.cluster
	}

	if len(quadKey) > iter {
		c := string(quadKey[iter])
		k, _ := strconv.Atoi(c)
		iter++
		if self.edges[k] != nil {
			return self.edges[k].getCluster(quadKey, iter)
		}
	}

	return nil
}
