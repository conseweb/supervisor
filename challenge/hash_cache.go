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
)

type BlocksHashCache interface {
	GetFromBlocksHashCache(highBlockNumber, lowBlockBumber uint64, hashType pb.HashType) (string, bool)
	SetBlocksHashToCache(highBlockNumber, lowBlockBumber uint64, hashType pb.HashType, hash string) bool
}

type defaultBlocksHashCache struct {
	caches map[string]*blocksHash
}

type blocksHash struct {
	blocksRange *pb.BlocksRange
	hashType    pb.HashType
	hash        string
}

func (this *defaultBlocksHashCache) GetFromBlocksHashCache(highBlockNumber, lowBlockBumber uint64, hashType pb.HashType) (string, bool) {
	key := this.blocksHashCacheKey(highBlockNumber, lowBlockBumber, hashType)

	cache, ok := this.caches[key]
	if !ok {
		return "", false
	}

	return cache.hash, true
}

func (this *defaultBlocksHashCache) SetBlocksHashToCache(highBlockNumber, lowBlockBumber uint64, hashType pb.HashType, hash string) bool {
	key := this.blocksHashCacheKey(highBlockNumber, lowBlockBumber, hashType)

	if _, ok := this.caches[key]; !ok {
		this.caches[key] = &blocksHash{
			blocksRange: &pb.BlocksRange{
				HighBlockNumber: highBlockNumber,
				LowBlockNumber:  lowBlockBumber,
			},
			hashType: hashType,
			hash:     hash,
		}

		return true
	}

	return false
}

func (this *defaultBlocksHashCache) blocksHashCacheKey(highBlockNumber, lowBlockBumber uint64, hashType pb.HashType) string {
	return HASH(pb.HashType_SHA256, []byte(fmt.Sprintf("%v/%v/%s", highBlockNumber, lowBlockBumber, hashType.String())))
}

func NewDefaultBlocksHashCache() BlocksHashCache {
	return &defaultBlocksHashCache{
		caches: make(map[string]*blocksHash),
	}
}

var blocksHashCache BlocksHashCache

func init() {
	blocksHashCache = NewDefaultBlocksHashCache()
}
