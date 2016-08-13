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

// once a farmer required to challenge the blocks hash,
// supervisor will store the random blocks range into cache,
// in order to avoid farmer fake requests
type FarmerChallengeCache interface {
	SetFarmerChallengeReq(farmerId string, highBlockNumber, lowBlockNumber uint64, hashAlgo pb.HashAlgo) bool
	GetFarmerChallengeReq(farmerId string, highBlockNumber, lowBlockNumber uint64, hashAlgo pb.HashAlgo) (*FarmerChallengeReq, bool)
	DelFarmerChallengeReq(farmerId string, highBlockNumber, lowBlockNumber uint64, hashAlgo pb.HashAlgo)
}

type defaultFarmerChallengeReqCache struct {
	caches map[string]*FarmerChallengeReq
}

type FarmerChallengeReq struct {
	farmerId    string
	blocksRange *pb.BlocksRange
	hashAlgo    pb.HashAlgo
}

func (this *defaultFarmerChallengeReqCache) cachekey(farmerId string, highBlockNumber, lowBlockNumber uint64, hashAlgo pb.HashAlgo) string {
	return HASH(pb.HashAlgo_SHA256, []byte(fmt.Sprintf("%s/%v/%v/%s", farmerId, highBlockNumber, lowBlockNumber, hashAlgo.String())))
}

func (this *defaultFarmerChallengeReqCache) SetFarmerChallengeReq(farmerId string, highBlockNumber, lowBlockNumber uint64, hashAlgo pb.HashAlgo) bool {
	key := this.cachekey(farmerId, highBlockNumber, lowBlockNumber, hashAlgo)

	if _, ok := this.caches[key]; !ok {
		logger.Debugf("challengeReq(%s) set to the cache", key)
		this.caches[key] = &FarmerChallengeReq{
			farmerId: farmerId,
			blocksRange: &pb.BlocksRange{
				HighBlockNumber: highBlockNumber,
				LowBlockNumber:  lowBlockNumber,
			},
			hashAlgo: hashAlgo,
		}

		return true
	}

	return false
}

func (this *defaultFarmerChallengeReqCache) GetFarmerChallengeReq(farmerId string, highBlockNumber, lowBlockNumber uint64, hashAlgo pb.HashAlgo) (*FarmerChallengeReq, bool) {
	key := this.cachekey(farmerId, highBlockNumber, lowBlockNumber, hashAlgo)

	cache, ok := this.caches[key]
	if !ok {
		logger.Debugf("challengeReq(%s) didn't hit the cache", key)
		return nil, false
	}

	logger.Debugf("challengeReq(%s) hit the cache: %v", key, cache)
	return cache, true
}

func (this *defaultFarmerChallengeReqCache) DelFarmerChallengeReq(farmerId string, highBlockNumber, lowBlockNumber uint64, hashAlgo pb.HashAlgo) {
	key := this.cachekey(farmerId, highBlockNumber, lowBlockNumber, hashAlgo)

	delete(this.caches, key)
}

func newDefaultFarmerChallengeReqCache() FarmerChallengeCache {
	return &defaultFarmerChallengeReqCache{
		caches: make(map[string]*FarmerChallengeReq),
	}
}

var (
	farmerChallengeReqCache FarmerChallengeCache
	challengeOnce           *sync.Once
)

func GetFarmerChallengeReqCache() FarmerChallengeCache {
	challengeOnce.Do(func() {
		farmerChallengeReqCache = newDefaultFarmerChallengeReqCache()
	})

	return farmerChallengeReqCache
}

func init() {
	challengeOnce = &sync.Once{}
}