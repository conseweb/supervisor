/*
Copyright Mojing Inc. 2016 All Rights Reserved.
Written by mint.zhao.chiu@gmail.com. github.com: https://www.github.com/mintzhao

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
package account

import (
	"errors"
)

var (
	errKeyNotFound = errors.New("key not found")
)

// key support ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789(62)
type AccountTreeNode struct {
	key      byte
	value    *FarmerAccountHandler
	parent   *AccountTreeNode
	depth    int
	children map[byte]*AccountTreeNode
	term     bool
}

func (this *AccountTreeNode) NewChild(key byte, handler *FarmerAccountHandler, term bool) *AccountTreeNode {
	node := &AccountTreeNode{
		key:      key,
		value:    handler,
		parent:   this,
		depth:    this.depth + 1,
		children: make(map[byte]*AccountTreeNode),
		term:     term,
	}

	this.children[key] = node
	return node
}

func (this *AccountTreeNode) DeleteChild(key byte) {
	delete(this.children, key)
}

type AccountTree struct {
	root *AccountTreeNode
	size int
}

const null = 0x0

// create a new trie tree with an initialized root node
func NewAccountTree() *AccountTree {
	return &AccountTree{
		root: &AccountTreeNode{
			children: make(map[byte]*AccountTreeNode),
		},
		size: 0,
	}
}

// put the key to the trie tree with value
func (t *AccountTree) Put(key string, handler *FarmerAccountHandler) error {
	if _, err := t.Get(key); err == errKeyNotFound {
		t.size++
	}

	bytesKeys := []byte(key)
	node := t.root
	for _, bkey := range bytesKeys {
		if n, ok := node.children[bkey]; ok {
			node = n
		} else {
			node = node.NewChild(bkey, handler, false)
		}
	}

	node = node.NewChild(null, handler, true)
	return nil
}

func (t *AccountTree) Get(key string) (*FarmerAccountHandler, error) {
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

func (t *AccountTree) Delete(key string) {
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

func (t *AccountTree) HasKeysWithPrefix(key string) bool {
	_, exist := getNode(t.root, []byte(key))
	return exist
}

func (t *AccountTree) PrefixSearch(pre string) []string {
	node, exist := getNode(t.root, []byte(pre))
	if !exist {
		return nil
	}

	return collect(node)
}

func (t *AccountTree) Keys() []string {
	return t.PrefixSearch("")
}

func (t *AccountTree) Len() int {
	return t.size
}

func getNode(node *AccountTreeNode, bkeys []byte) (*AccountTreeNode, bool) {
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

func collect(node *AccountTreeNode) []string {
	keys := []string{}
	i := int(0)

	nodes := []*AccountTreeNode{node}
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
				word = string(p.key) + word
			}

			keys = append(keys, word)
		}
	}

	return keys
}
