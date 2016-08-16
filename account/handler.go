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
	"github.com/conseweb/supervisor/challenge"
	pb "github.com/conseweb/supervisor/protos"
	"github.com/golang/protobuf/proto"
	"github.com/looplab/fsm"
	"github.com/spf13/viper"
	"math/rand"
	"time"
)

type FarmerAccountHandler struct {
	account      *pb.FarmerAccount
	fsm          *fsm.FSM
	lostCount    int
	nextPingTime int64
}

// whether or not need challenge blocks hash
func (this *FarmerAccountHandler) needChallengeBlocks(highBlockNumber, lowBlockBumber uint64) (need bool, brange *pb.BlocksRange) {
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
func (this *FarmerAccountHandler) challengeHashAlgo() pb.HashAlgo {
	hashAlgo := viper.GetString("farmer.challenge.hashalgo")
	if hashAlgo == "" {
		hashAlgo = pb.HashAlgo_SHA256.String()
	}
	return pb.HashAlgo(pb.HashAlgo_value[hashAlgo])
}

// randomly return next ping time
func (this *FarmerAccountHandler) NextPingTime() int64 {
	if this.nextPingTime <= 0 {
		this.nextPingTime = time.Now().Add(nextPingInterval()).UnixNano()
	}

	return this.nextPingTime
}

func (this *FarmerAccountHandler) Account() *pb.FarmerAccount {
	if this.account == nil {
		return nil
	}
	return this.account
}

// after online, we set farmer's lost count 0
func (this *FarmerAccountHandler) afterEvent() {
	this.account.LastModifiedTime = time.Now().UnixNano()
	this.account.State = pb.FarmerState(pb.FarmerState_value[this.fsm.Current()])

	time.AfterFunc(nextPingInterval(), func() {
		this.lostCount++
		if this.fsm.Can("lost") {
			this.Lost()
		}
		if this.lostCount >= viper.GetInt("farmer.ping.lostcount") {
			if this.fsm.Can("offline") {
				this.OffLine()
			}
		}
	})

	UpdateFarmerHandler(this)
}

func (this *FarmerAccountHandler) OnLine() error {
	if err := this.fsm.Event("online"); err != nil {
		logger.Errorf("farmer online return err: %v", err)
		return err
	}

	this.lostCount = 0
	this.afterEvent()

	return nil
}

func (this *FarmerAccountHandler) Ping(highBlockNumber, lowBlockNumber uint64) (need bool, brange *pb.BlocksRange, hashAlgo pb.HashAlgo, err error) {
	if this.fsm.Current() == pb.FarmerState_OFFLINE.String() {
		err = errors.New("farmer is offline")
		return
	}

	need, brange = this.needChallengeBlocks(highBlockNumber, lowBlockNumber)
	if need {
		hashAlgo = this.challengeHashAlgo()
	} else {
		// if no need to challenge, just add balance this time
		this.calcBalance()
	}

	this.lostCount = 0
	this.afterEvent()

	return nil
}

func (this *FarmerAccountHandler) ConquerChallenge(highBlockNumber, lowBlockNumber uint64, hashAlgo pb.HashAlgo, blocksHash string) error {
	if challenge.ConquerChallenge(this.account.FarmerID, highBlockNumber, lowBlockNumber, hashAlgo, blocksHash) {
		// TODO calc farmer balance
		this.calcBalance()
	} else {
		this.account.Balance = 0
		return errors.New("farmer conquer challenge fail")
	}

	this.lostCount = 0
	this.account.LastChallengeTime = time.Now().UnixNano()
	this.afterEvent()

	return nil
}

func (this *FarmerAccountHandler) calcBalance() {
	this.account.Balance += 100
}

func (this *FarmerAccountHandler) Lost() error {
	if this.lostCount <= 0 {
		return errors.New("current lost count <= 0")
	}

	if err := this.fsm.Event("lost"); err != nil {
		logger.Errorf("farmer lost return err: %v", err)
		return err
	}

	this.afterEvent()

	return nil
}

func (this *FarmerAccountHandler) OffLine() error {
	if err := this.fsm.Event("offline"); err != nil {
		logger.Errorf("farmer offline return err: %v", err)
		return err
	}

	this.lostCount = 0
	this.afterEvent()

	return nil
}

func (this *FarmerAccountHandler) beforeEvent(e *fsm.Event) {
	logger.Debugf("farmer(%v) do event(%v) from state %v enter state %v", this.account.FarmerID, e.Event, e.Src, e.Dst)
}

// TODO change farmerid 2 backend uniquekey
// at current farmerid may be a short string/int
// this func translate farmerid 2 backend unique key
// now just return farmerid
func farmerId2Key(farmerId string) string {
	return farmerId
}

// marshal farmer account 2 proto bytes
// if there is an error, return nil
func (this *FarmerAccountHandler) marshal() ([]byte, error) {
	accountBytes, err := proto.Marshal(this.account)
	if err != nil {
		logger.Errorf("marshal farmer account err: %v", err)
		return nil, err
	}

	return accountBytes, nil
}

// unmarshal proto bytes 2 farmer account
func (this *FarmerAccountHandler) unmarshal(fBytes []byte) error {
	this.account = &pb.FarmerAccount{}
	if err := proto.Unmarshal(fBytes, this.account); err != nil {
		logger.Errorf("unmarshal farmer account err: %v", err)
		return err
	}

	return nil
}

func nextPingInterval() time.Duration {
	basicInterval := viper.GetInt("farmer.ping.interval")
	if basicInterval == 0 {
		basicInterval = 900
	}

	return time.Duration(basicInterval) * time.Second

	//up := viper.GetInt("farmer.ping.up")
	//down := viper.GetInt("farmer.ping.down")
	//
	//interval := basicInterval
	//upflag := (rand.Int() % 2) == 0
	//if upflag {
	//	interval += rand.Intn(up)
	//} else {
	//	interval -= rand.Intn(down)
	//}
	//
	//if interval < 0 {
	//	interval = basicInterval
	//}
	//
	//return time.Duration(interval) * time.Second
}
