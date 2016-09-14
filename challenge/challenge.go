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
package challenge

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	pb "github.com/conseweb/common/protos"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/op/go-logging"
)

var (
	logger         = logging.MustGetLogger("supervisor")
	ErrOutOfBounds = errors.New("supervisor/challenge: blocks out of bounds")
)

func ConquerChallenge(farmerId string, highBlockNumber, lowBlockNumber uint64, hashAlgo pb.HashAlgo, blocksHash string) bool {
	// get challenge request hash, if not found, mean there is no such challenge request, farmer fake it
	_, get := GetFarmerChallengeReqCache().GetFarmerChallengeReq(farmerId, highBlockNumber, lowBlockNumber, hashAlgo)
	if !get {
		logger.Errorf("supervisor/challenge: invalid challenge request. farmerId' %s, highBlockNumber: %d, lowBlockNumber: %d, hashAlgo: %v, blocksHash: %s", farmerId, highBlockNumber, lowBlockNumber, hashAlgo, blocksHash)
		return false
	}

	// once get challenge request from the cache, delete it, one request just can be fetch one time
	// TODO whether or not just move del into get
	GetFarmerChallengeReqCache().DelFarmerChallengeReq(farmerId, highBlockNumber, lowBlockNumber, hashAlgo)

	// get blocks hash from blocks hash cache, if not found, just hash it and put it into cache
	originalHash, get := GetBlocksHashCache().GetFromBlocksHashCache(highBlockNumber, lowBlockNumber, hashAlgo)
	if !get {
		blocksBytes, err := GetBlocksBytes(highBlockNumber, lowBlockNumber)
		if err != nil {
			logger.Errorf("get blocks[%d, %d] bytes err: %v", highBlockNumber, lowBlockNumber, err)
			return false
		}
		originalHash = HASH(hashAlgo, blocksBytes)
		GetBlocksHashCache().SetBlocksHashToCache(highBlockNumber, lowBlockNumber, hashAlgo, originalHash)
	}

	// compare farmer result & sv result
	serverHash := FarmerBindConquerHash(farmerId, hashAlgo, originalHash)
	if strings.Compare(blocksHash, serverHash) != 0 {
		logger.Warningf("farmer[%s] conquer challenge fail, farmer hash: %s, server hash: %s", farmerId, blocksHash, serverHash)
		return false
	}
	logger.Debugf("farmer[%s] conquer challenge success.", farmerId)

	return true
}

// compare hash isn't hash of blocks, is hash of (farmerId string and hash of blocks's string)
func FarmerBindConquerHash(farmerId string, hashAlgo pb.HashAlgo, originalHash string) string {
	return HASH(hashAlgo, bytes.NewBufferString(fmt.Sprintf("%s%s", originalHash, farmerId)).Bytes())
}

func GetBlocksBytes(highBlockNumber, lowBlockNumber uint64) ([]byte, error) {
	return []byte(""), nil

	blocksBuffer := bytes.NewBufferString("")
	ldg, err := ledger.GetLedger()
	if err != nil {
		return blocksBuffer.Bytes(), err
	}

	if highBlockNumber >= ldg.GetBlockchainSize() {
		return blocksBuffer.Bytes(), ErrOutOfBounds
	}
	if highBlockNumber < lowBlockNumber {
		return blocksBuffer.Bytes(), ErrOutOfBounds
	}

	currentBlock, err := ldg.GetBlockByNumber(highBlockNumber)
	if err != nil {
		return blocksBuffer.Bytes(), fmt.Errorf("Error fetching block %d.", highBlockNumber)
	}
	if currentBlock == nil {
		return blocksBuffer.Bytes(), fmt.Errorf("Block %d is nil.", highBlockNumber)
	}

	for i := highBlockNumber; i > lowBlockNumber; i-- {
		previousBlock, err := ldg.GetBlockByNumber(i - 1)
		if err != nil {
			return blocksBuffer.Bytes(), err
		}
		if previousBlock == nil {
			return blocksBuffer.Bytes(), fmt.Errorf("Block %d is nil.", i)
		}
		previousBlockHash, err := previousBlock.GetHash()
		if err != nil {
			return blocksBuffer.Bytes(), err
		}
		if bytes.Compare(previousBlockHash, currentBlock.PreviousBlockHash) != 0 {
			return blocksBuffer.Bytes(), fmt.Errorf("Blocks hash can not match.")
		}

		if currentBlockBytes, err := currentBlock.Bytes(); err != nil {
			return blocksBuffer.Bytes(), err
		} else {
			blocksBuffer.Write(currentBlockBytes)
		}
		currentBlock = previousBlock
	}

	return blocksBuffer.Bytes(), nil
}
