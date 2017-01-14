package trie

import (
	"errors"
	"strconv"
)

type IData interface{}

type Trie struct {
	root   *node
	length int
}

type node struct {
	NodeData
	edges    [4]*node
	prefixes int
}

type NodeData struct {
	word string
	data []IData
}

func NewTrie() *Trie {
	return &Trie{
		newNode(),
		0,
	}
}

func (self *Trie) AddWord(word string, data IData) error {
	if len(word) == 0 {
		return errors.New("empty word")
	}

	self.root.addWord(word, 0, data)
	return nil
}

func (self *Trie) Lookup(prefix string) ([]*NodeData, error) {
	var data []*NodeData
	if len(prefix) == 0 {
		return data, errors.New("empty prefix")
	}

	self.root.lookup(prefix, 0, &data)
	return data, nil
}

func (self *Trie) GetPrefixesCount(word string) int {
	return self.root.getPrefixesCount(word, 0)
}

func (self *NodeData) GetWord() string {
	return self.word
}

func (self *NodeData) GetData() []IData {
	return self.data
}

func newNode() *node {
	return &node{
		NodeData{"", make([]IData, 0)}, [4]*node{nil, nil, nil, nil}, 0,
	}
}

func (self *node) addWord(word string, iter int, data IData) {
	if len(word) <= iter {
		self.word = word
		self.data = append(self.data, data)
		return
	}

	self.prefixes++
	c := string(word[iter])
	k, _ := strconv.Atoi(c)
	if self.edges[k] == nil {
		self.edges[k] = newNode()
	}

	iter++
	self.edges[k].addWord(word, iter, data)
}

func (self *node) lookup(prefix string, iter int, data *[]*NodeData) {
	if len(self.word) > 0 && len(prefix) == iter {
		*data = append(*data, &NodeData{self.word, self.data})
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

func (self *node) getPrefixesCount(word string, iter int) int {
	if len(word) == iter {
		return self.prefixes
	}

	if len(word) > iter {
		c := string(word[iter])
		k, _ := strconv.Atoi(c)
		iter++
		if self.edges[k] != nil {
			return self.edges[k].getPrefixesCount(word, iter)
		}
	}

	return 0
}
