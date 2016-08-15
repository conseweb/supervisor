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
	"fmt"
	pb "github.com/conseweb/supervisor/protos"
	"sync"
)

type BlocksHashCache interface {
	GetFromBlocksHashCache(highBlockNumber, lowBlockBumber uint64, hashAlgo pb.HashAlgo) (string, bool)
	SetBlocksHashToCache(highBlockNumber, lowBlockBumber uint64, hashAlgo pb.HashAlgo, hash string) bool
	Close() error
}

type defaultBlocksHashCache struct {
	caches map[string]*blocksHashItem
}

type blocksHashItem struct {
	blocksRange *pb.BlocksRange
	hashAlgo    pb.HashAlgo
	hash        string
}

func (this *defaultBlocksHashCache) GetFromBlocksHashCache(highBlockNumber, lowBlockBumber uint64, hashAlgo pb.HashAlgo) (string, bool) {
	key := this.blocksHashCacheKey(highBlockNumber, lowBlockBumber, hashAlgo)

	cache, ok := this.caches[key]
	if !ok {
		logger.Debugf("blockshash(%s) didn't hit the cache", key)
		return "", false
	}

	logger.Debugf("blockshash(%s) hit the cache: %v", key, cache)
	return cache.hash, true
}

func (this *defaultBlocksHashCache) SetBlocksHashToCache(highBlockNumber, lowBlockBumber uint64, hashAlgo pb.HashAlgo, hash string) bool {
	key := this.blocksHashCacheKey(highBlockNumber, lowBlockBumber, hashAlgo)

	if _, ok := this.caches[key]; !ok {
		logger.Debugf("blockshash(%s) set to the cache", key)
		this.caches[key] = &blocksHashItem{
			blocksRange: &pb.BlocksRange{
				HighBlockNumber: highBlockNumber,
				LowBlockNumber:  lowBlockBumber,
			},
			hashAlgo: hashAlgo,
			hash:     hash,
		}

		return true
	}

	return false
}

func (this *defaultBlocksHashCache) Close() error {
	this.caches = nil
	return nil
}

func (this *defaultBlocksHashCache) blocksHashCacheKey(highBlockNumber, lowBlockBumber uint64, hashAlgo pb.HashAlgo) string {
	return HASH(pb.HashAlgo_SHA256, []byte(fmt.Sprintf("%v/%v/%s", highBlockNumber, lowBlockBumber, hashAlgo.String())))
}

var (
	blocksHashCache BlocksHashCache
	hashonce        *sync.Once
)

func newDefaultBlocksHashCache() BlocksHashCache {
	return &defaultBlocksHashCache{
		caches: make(map[string]*blocksHashItem),
	}
}

func GetBlocksHashCache() BlocksHashCache {
	hashonce.Do(func() {
		blocksHashCache = newDefaultBlocksHashCache()
	})

	return blocksHashCache
}

func init() {
	hashonce = &sync.Once{}
}
