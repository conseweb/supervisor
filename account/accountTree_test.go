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
	"fmt"
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type TestAccountTree struct{}

var _ = check.Suite(&TestAccountTree{})

func (s *TestAccountTree) TestTrieTree_Put(c *check.C) {
	trie := NewAccountTree()
	c.Check(trie.Put("abc", nil), check.IsNil)
}

func (s *TestAccountTree) TestTrieTree_Get(c *check.C) {
	trie := NewAccountTree()

	c.Check(trie.Put("abc", &FarmerAccountHandler{}), check.IsNil)
	vget, err := trie.Get("abc")
	c.Check(err, check.IsNil)
	c.Assert(vget, check.NotNil)
}

func (s *TestAccountTree) TestTrieTree_Delete(c *check.C) {
	trie := NewAccountTree()
	initial := []string{"football", "foostar", "foosball"}

	for _, key := range initial {
		c.Check(trie.Put(key, nil), check.IsNil)
	}
	trie.Delete("football")

	keys := trie.Keys()
	c.Assert(len(keys), check.Equals, 2)
}

func (s *TestAccountTree) TestTrieTree_HasKeysWithPrefix(c *check.C) {
	trie := NewAccountTree()

	c.Check(trie.Put("abc", nil), check.IsNil)
	c.Check(trie.Put("acb", nil), check.IsNil)
	c.Check(trie.HasKeysWithPrefix("ab"), check.Equals, true)
	c.Check(trie.HasKeysWithPrefix("ac"), check.Equals, true)
	c.Check(trie.HasKeysWithPrefix("bc"), check.Equals, false)
}

func (s *TestAccountTree) TestTrieTree_PrefixSearch(c *check.C) {
	trie := NewAccountTree()

	c.Check(trie.Put("abc", nil), check.IsNil)
	c.Check(trie.Put("acb", nil), check.IsNil)
	c.Check(2, check.Equals, len(trie.PrefixSearch("a")))
	c.Check(2, check.Not(check.Equals), len(trie.PrefixSearch("b")))
}

func (s *TestAccountTree) TestTrieTree_Update(c *check.C) {
	trie := NewAccountTree()

	c.Check(trie.Put("abc", nil), check.IsNil)
	get1, err1 := trie.Get("abc")
	c.Check(err1, check.IsNil)
	c.Check(get1, check.IsNil)

	c.Check(trie.Put("abc", &FarmerAccountHandler{}), check.IsNil)
	get2, err2 := trie.Get("abc")
	c.Check(err2, check.IsNil)
	c.Check(get2, check.NotNil)
}

func (s *TestAccountTree) TestTrieTree_Len(c *check.C) {
	trie := NewAccountTree()

	c.Check(trie.Put("abc", nil), check.IsNil)
	c.Check(trie.Put("abc", nil), check.IsNil)
	c.Check(trie.Put("abd", nil), check.IsNil)

	c.Assert(trie.Len(), check.Equals, 2)
}

func (s *TestAccountTree) BenchmarkTrieTree_Put(c *check.C) {
	trie := NewAccountTree()
	for i := 0; i < c.N; i++ {
		trie.Put(fmt.Sprintf("abc%v", i), nil)
	}
}

func (s *TestAccountTree) BenchmarkTrieTree_Get(c *check.C) {
	trie := NewAccountTree()
	trie.Put("abc", nil)
	for i := 0; i < c.N; i++ {
		trie.Get("abc")
	}
}

func (s *TestAccountTree) BenchmarkTrieTree_Delete(c *check.C) {
	trie := NewAccountTree()
	for i := 0; i < c.N; i++ {
		trie.Put("abcdefg", nil)
		trie.Delete("abcdefg")
	}
}
