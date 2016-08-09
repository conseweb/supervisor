/*
Copyright Mojing Inc. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package tree

import (
	"errors"
)

var (
	errKeyNotFound = errors.New("key not found")
)

// key support ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789(62)
type TrieTreeNode struct {
	bkey     byte
	value    []byte
	parent   *TrieTreeNode
	depth    int
	children map[byte]*TrieTreeNode
	term     bool
}

func (n *TrieTreeNode) NewChild(bkey byte, value []byte, term bool) *TrieTreeNode {
	node := &TrieTreeNode{
		bkey:     bkey,
		value:    value,
		parent:   n,
		depth:    n.depth + 1,
		children: make(map[byte]*TrieTreeNode),
		term:     term,
	}

	n.children[bkey] = node
	return node
}

func (n *TrieTreeNode) DeleteChild(bkey byte) {
	delete(n.children, bkey)
}

type TrieTree struct {
	root *TrieTreeNode
	size int
}

const null = 0x0

// create a new trie tree with an initialized root node
// and must provide a marshalFunc to marshal value into node
func NewTrie() *TrieTree {
	return &TrieTree{
		root: &TrieTreeNode{
			children: make(map[byte]*TrieTreeNode),
		},
		size: 0,
	}
}

// put the key to the trie tree with value
func (t *TrieTree) Put(key string, value []byte) error {
	if _, err := t.Get(key); err == errKeyNotFound {
		t.size++
	}

	bytesKeys := []byte(key)
	node := t.root
	for _, bkey := range bytesKeys {
		if n, ok := node.children[bkey]; ok {
			node = n
		} else {
			node = node.NewChild(bkey, value, false)
		}
	}

	node = node.NewChild(null, value, true)
	return nil
}

func (t *TrieTree) Get(key string) ([]byte, error) {
	node, geted := getNode(t.root, []byte(key))
	if !geted {
		return nil, errKeyNotFound
	}

	node, ok := node.children[null]
	if !ok || !node.term {
		return nil, errKeyNotFound
	}

	return node.value, nil
}

func (t *TrieTree) Delete(key string) {
	i := int(0)
	bkeys := []byte(key)
	node, exist := getNode(t.root, bkeys)
	if !exist {
		return
	}

	t.size--
	for n := node.parent; n != nil; n = n.parent {
		i++
		if len(n.children) > 1 {
			n.DeleteChild(bkeys[len(bkeys)-i])
			break
		}
	}
}

func (t *TrieTree) HasKeysWithPrefix(key string) bool {
	_, exist := getNode(t.root, []byte(key))
	return exist
}

func (t *TrieTree) PrefixSearch(pre string) []string {
	node, exist := getNode(t.root, []byte(pre))
	if !exist {
		return nil
	}

	return collect(node)
}

func (t *TrieTree) Keys() []string {
	return t.PrefixSearch("")
}

func (t *TrieTree) Len() int {
	return t.size
}

func getNode(node *TrieTreeNode, bkeys []byte) (*TrieTreeNode, bool) {
	if node == nil {
		return nil, false
	}

	if len(bkeys) == 0 {
		return node, true
	}

	n, ok := node.children[bkeys[0]]
	if !ok {
		return nil, false
	}

	var nbkeys []byte
	if len(bkeys) > 1 {
		nbkeys = bkeys[1:]
	} else {
		nbkeys = bkeys[0:0]
	}

	return getNode(n, nbkeys)
}

func collect(node *TrieTreeNode) []string {
	keys := []string{}
	i := int(0)

	nodes := []*TrieTreeNode{node}
	for l := len(nodes); l != 0; l = len(nodes) {
		i = l - 1
		n := nodes[i]
		nodes = nodes[:i]
		for _, cn := range n.children {
			nodes = append(nodes, cn)
		}

		if n.term {
			word := ""
			for p := n.parent; p.depth != 0; p = p.parent {
				word = string(p.bkey) + word
			}

			keys = append(keys, word)
		}
	}

	return keys
}

