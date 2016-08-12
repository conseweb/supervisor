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
	"fmt"
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type TestTrie struct{}

var _ = check.Suite(&TestTrie{})

func (s *TestTrie) TestTrieTree_Put(c *check.C) {
	trie := NewTrie()
	c.Check(trie.Put("abc", nil), check.IsNil)
}

func (s *TestTrie) TestTrieTree_Get(c *check.C) {
	trie := NewTrie()

	c.Check(trie.Put("abc", []byte("abc")), check.IsNil)
	vget, err := trie.Get("abc")
	c.Check(err, check.IsNil)
	c.Assert(string(vget), check.Equals, "abc")
}

func (s *TestTrie) TestTrieTree_Delete(c *check.C) {
	trie := NewTrie()
	initial := []string{"football", "foostar", "foosball"}

	for _, key := range initial {
		c.Check(trie.Put(key, nil), check.IsNil)
	}
	trie.Delete("football")

	keys := trie.Keys()
	c.Assert(len(keys), check.Equals, 2)
}

func (s *TestTrie) TestTrieTree_HasKeysWithPrefix(c *check.C) {
	trie := NewTrie()

	c.Check(trie.Put("abc", nil), check.IsNil)
	c.Check(trie.Put("acb", nil), check.IsNil)
	c.Check(trie.HasKeysWithPrefix("ab"), check.Equals, true)
	c.Check(trie.HasKeysWithPrefix("ac"), check.Equals, true)
	c.Check(trie.HasKeysWithPrefix("bc"), check.Equals, false)
}

func (s *TestTrie) TestTrieTree_PrefixSearch(c *check.C) {
	trie := NewTrie()

	c.Check(trie.Put("abc", nil), check.IsNil)
	c.Check(trie.Put("acb", nil), check.IsNil)
	c.Check(2, check.Equals, len(trie.PrefixSearch("a")))
	c.Check(2, check.Not(check.Equals), len(trie.PrefixSearch("b")))
}

func (s *TestTrie) TestTrieTree_Update(c *check.C) {
	trie := NewTrie()

	c.Check(trie.Put("abc", []byte("abc")), check.IsNil)
	get1, err1 := trie.Get("abc")
	c.Check(err1, check.IsNil)
	c.Check(string(get1), check.Equals, "abc")

	c.Check(trie.Put("abc", []byte("acb")), check.IsNil)
	get2, err2 := trie.Get("abc")
	c.Check(err2, check.IsNil)
	c.Check(string(get2), check.Equals, "acb")
	c.Check(string(get2), check.Not(check.Equals), "abc")
}

func (s *TestTrie) TestTrieTree_Len(c *check.C) {
	trie := NewTrie()

	c.Check(trie.Put("abc", []byte("abc")), check.IsNil)
	c.Check(trie.Put("abc", []byte("abc")), check.IsNil)
	c.Check(trie.Put("abd", []byte("abc")), check.IsNil)

	c.Assert(trie.Len(), check.Equals, 2)
}

func (s *TestTrie) BenchmarkTrieTree_Put(c *check.C) {
	trie := NewTrie()
	for i := 0; i < c.N; i++ {
		trie.Put(fmt.Sprintf("abcdefghijklmnopqrstuvwxyz%v", i), []byte("abc"))
	}
}

func (s *TestTrie) BenchmarkTrieTree_Get(c *check.C) {
	trie := NewTrie()
	trie.Put("abcdefg", nil)
	for i := 0; i < c.N; i++ {
		trie.Get("abcdefg")
	}
}

func (s *TestTrie) BenchmarkTrieTree_Delete(c *check.C) {
	trie := NewTrie()
	for i := 0; i < c.N; i++ {
		trie.Put("abcdefg", nil)
		trie.Delete("abcdefg")
	}
}
