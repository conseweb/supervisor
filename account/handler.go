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
package account

import (
	"errors"
	"math/rand"
	"time"

	"github.com/conseweb/supervisor/challenge"
	pb "github.com/conseweb/supervisor/protos"
	"github.com/golang/protobuf/proto"
	"github.com/looplab/fsm"
	"github.com/spf13/viper"
)

type FarmerAccountHandler struct {
	account                *pb.FarmerAccount
	fsm                    *fsm.FSM
	lostCount              int
	nextPingTime           int64
	nextConquerTime        int64
	nextFarmerChallengeReq *challenge.FarmerChallengeReq
}

// whether or not need challenge blocks hash
func (h *FarmerAccountHandler) needChallengeBlocks(highBlockNumber, lowBlockBumber uint64) (need bool, brange *pb.BlocksRange) {
	brange = &pb.BlocksRange{}

	need = (rand.Int() % 2) == 0
	if need {
		if highBlockNumber < lowBlockBumber {
			return
		}

		brange.HighBlockNumber = highBlockNumber - uint64(rand.Int63n(int64(highBlockNumber-lowBlockBumber)))
		brange.LowBlockNumber = lowBlockBumber + uint64(rand.Int63n(int64(brange.HighBlockNumber-lowBlockBumber)))
	}

	return
}

// choose challenge hash type from viper
func (h *FarmerAccountHandler) challengeHashAlgo() pb.HashAlgo {
	hashAlgo := viper.GetString("farmer.challenge.hashalgo")
	if hashAlgo == "" {
		hashAlgo = pb.HashAlgo_SHA256.String()
	}
	return pb.HashAlgo(pb.HashAlgo_value[hashAlgo])
}

// randomly return next ping time
func (h *FarmerAccountHandler) NextPingTime() int64 {
	if h.nextPingTime <= 0 {
		h.nextPingTime = time.Now().Add(nextPingInterval()).UnixNano()
	}

	return h.nextPingTime
}

func (h *FarmerAccountHandler) Account() *pb.FarmerAccount {
	if h.account == nil {
		return nil
	}
	return h.account
}

// after online, we set farmer's lost count 0
func (h *FarmerAccountHandler) afterEvent() {
	h.account.LastModifiedTime = time.Now().UnixNano()
	h.account.State = pb.FarmerState(pb.FarmerState_value[h.fsm.Current()])

	UpdateFarmerHandler(h)
}

func (h *FarmerAccountHandler) OnLine() error {
	if h.fsm.Can("online") {
		if err := h.fsm.Event("online"); err != nil {
			logger.Errorf("farmer online return err: %v", err)
			return err
		}
	} else {
		return errors.New("already online, can not override.")
	}

	h.lostCount = 0
	h.afterEvent()

	return nil
}

func (h *FarmerAccountHandler) Ping(highBlockNumber, lowBlockNumber uint64) (need bool, brange *pb.BlocksRange, hashAlgo pb.HashAlgo, err error) {
	if h.fsm.Current() == pb.FarmerState_OFFLINE.String() {
		err = errors.New("farmer is offline")
		return
	}

	need, brange = h.needChallengeBlocks(highBlockNumber, lowBlockNumber)
	if need {
		hashAlgo = h.challengeHashAlgo()

		// sv cache challenge req
		req, set := challenge.GetFarmerChallengeReqCache().SetFarmerChallengeReq(h.account.FarmerID, brange.HighBlockNumber, brange.LowBlockNumber, hashAlgo)
		if set {
			// set handler's nextConquerTime and nextChallengeReq
			h.nextConquerTime = time.Now().Add(nextChallengeDelay()).UnixNano()
			h.nextFarmerChallengeReq = req
		}
	} else {
		// if no need to challenge, just add balance h time
		h.calcBalance()
	}

	h.lostCount = 0
	h.afterEvent()

	return
}

func (h *FarmerAccountHandler) ConquerChallenge(highBlockNumber, lowBlockNumber uint64, hashAlgo pb.HashAlgo, blocksHash string) error {
	if challenge.ConquerChallenge(h.account.FarmerID, highBlockNumber, lowBlockNumber, hashAlgo, blocksHash) {
		h.calcBalance()
	} else {
		h.punishBalance()
		return errors.New("farmer conquer challenge fail")
	}

	h.lostCount = 0
	h.account.LastChallengeTime = time.Now().UnixNano()
	h.nextConquerTime = 0
	h.nextFarmerChallengeReq = nil
	h.afterEvent()

	return nil
}

// TODO calc farmer balance
func (h *FarmerAccountHandler) calcBalance() {
	h.account.Balance += 100
}

func (h *FarmerAccountHandler) punishBalance() {
	h.account.Balance = 0
}

func (h *FarmerAccountHandler) Lost() error {
	if h.lostCount <= 0 {
		return errors.New("current lost count <= 0")
	}

	if h.fsm.Can("lost") {
		if err := h.fsm.Event("lost"); err != nil {
			logger.Errorf("farmer lost return err: %v", err)
			return err
		}
	}

	h.afterEvent()

	return nil
}

func (h *FarmerAccountHandler) OffLine() error {
	if h.fsm.Can("offline") {
		if err := h.fsm.Event("offline"); err != nil {
			logger.Errorf("farmer offline return err: %v", err)
			return err
		}
	}

	h.afterEvent()

	return nil
}

func (h *FarmerAccountHandler) beforeEvent(e *fsm.Event) {
	logger.Debugf("farmer(%v) do event(%v) from state %v enter state %v", h.account.FarmerID, e.Event, e.Src, e.Dst)
}

// TODO change farmerid 2 backend uniquekey
// at current farmerid may be a short string/int
// h func translate farmerid 2 backend unique key
// now just return farmerid
func farmerId2Key(farmerId string) string {
	return farmerId
}

// marshal farmer account 2 proto bytes
func farmerAccount2Bytes(account *pb.FarmerAccount) ([]byte, error) {
	accountBytes, err := proto.Marshal(account)
	if err != nil {
		logger.Errorf("marshal farmer account err: %v", err)
		return nil, err
	}

	return accountBytes, nil
}

// unmarshal proto bytes 2 farmer account
func bytes2FarmerAccount(fBytes []byte) (*pb.FarmerAccount, error) {
	account := &pb.FarmerAccount{}
	if err := proto.Unmarshal(fBytes, account); err != nil {
		logger.Errorf("unmarshal farmer account err: %v", err)
		return nil, err
	}

	return account, nil
}

func nextPingInterval() time.Duration {
	if interval, err := time.ParseDuration(viper.GetString("farmer.ping.interval")); err == nil {
		return interval
	}

	viper.Set("farmer.ping.interval", "900s")
	return time.Duration(900) * time.Second
}

func nextChallengeDelay() time.Duration {
	if delay, err := time.ParseDuration(viper.GetString("farmer.challenge.delay")); err == nil {
		return delay
	}

	viper.Set("farmer.challenge.delay", "10s")
	return time.Duration(10) * time.Second
}
