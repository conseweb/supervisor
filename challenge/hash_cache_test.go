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
package challenge

import (
	pb "github.com/conseweb/supervisor/protos"
	"gopkg.in/check.v1"
)

type BlocksHashCacheTest struct {
	cache *defaultBlocksHashCache
}

var _ = check.Suite(&BlocksHashCacheTest{})

func (this *BlocksHashCacheTest) TestDefaultBlocksHashCacheSet(c *check.C) {
	cache := NewDefaultBlocksHashCache()

	c.Check(cache.SetBlocksHashToCache(100, 20, pb.HashType_SHA1, "pretend as hash"), check.Equals, true)
	c.Check(cache.SetBlocksHashToCache(100, 20, pb.HashType_SHA1, "pretend as hash"), check.Equals, false)
}

func (this *BlocksHashCacheTest) TestDefaultBlocksHashCacheGet(c *check.C) {
	cache := NewDefaultBlocksHashCache()
	cache.SetBlocksHashToCache(100, 20, pb.HashType_SHA1, "pretend as hash")

	hash, getted := cache.GetFromBlocksHashCache(100, 20, pb.HashType_SHA1)
	c.Check(getted, check.Equals, true)
	c.Check(hash, check.Equals, "pretend as hash")
}

func (this *BlocksHashCacheTest) BenchmarkDefaultBlocksHashCacheSet(c *check.C) {
	cache := NewDefaultBlocksHashCache()
	for i := 0; i < c.N; i++ {
		cache.SetBlocksHashToCache(100+uint64(i), 20, pb.HashType_SHA1, "pretend as hash")
	}
}

func (this *BlocksHashCacheTest) BenchmarkDefaultBlocksHashCacheGet(c *check.C) {
	cache := NewDefaultBlocksHashCache()
	for i := 0; i < c.N; i++ {
		cache.GetFromBlocksHashCache(100+uint64(i), 20, pb.HashType_SHA1)
	}
}
