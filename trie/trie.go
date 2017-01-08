package trie

import (
	"errors"
	"strconv"
)

type Trie struct {
	root   *node
	length int
}

type node struct {
	word  string
	edges [4]*node
}

func NewTrie() *Trie {
	return &Trie{
		newNode(),
		0,
	}
}

func (self *Trie) AddWord(word string) error {
	if len(word) == 0 {
		return errors.New("empty word")
	}

	self.root.addWord(word, 0)
	return nil
}

func (self *Trie) Lookup(prefix string) ([]string, error) {
	prefixes := make([]string, 0)
	if len(prefix) == 0 {
		return prefixes, errors.New("empty prefix")
	}

	self.root.lookup(prefix, &prefixes)
	return prefixes, nil
}

func newNode() *node {
	return &node{
		"", [4]*node{nil, nil, nil, nil},
	}
}

func (self *node) addWord(word string, iter int) {
	if len(word) <= iter {
		self.word = word
		return
	}

	c := string(word[iter])
	k, _ := strconv.Atoi(c)
	if self.edges[k] == nil {
		self.edges[k] = newNode()
	}

	iter++
	self.edges[k].addWord(word, iter)
}

func (self *node) lookup(prefix string, prefixes *[]string) {
	if len(self.word) > 0 && len(prefix) == 0 {
		*prefixes = append(*prefixes, self.word)
	}

	if len(prefix) != 0 {
		c, tail := string(prefix[0]), prefix[1:]
		k, _ := strconv.Atoi(c)
		if self.edges[k] != nil {
			self.edges[k].lookup(tail, prefixes)
		}

		return
	}

	for _, edge := range self.edges {
		if edge != nil {
			edge.lookup(prefix, prefixes)
		}
	}
}
