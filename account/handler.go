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
func (this *FarmerAccountHandler) NeedChallengeBlocks(highBlockNumber, lowBlockBumber uint64) (need bool, brange *pb.BlocksRange) {
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
func (this *FarmerAccountHandler) ChallengeHashType() pb.HashType {
	hashType := viper.GetString("farmer.challenge.hash")
	if hashType == "" {
		hashType = pb.HashType_SHA256.String()
	}
	return pb.HashType(pb.HashType_value[hashType])
}

// randomly return next ping time
func (this *FarmerAccountHandler) NextPingTime() int64 {
	return this.nextPingTime
}

func (this *FarmerAccountHandler) Account() *pb.FarmerAccount {
	return this.account
}

func (this *FarmerAccountHandler) OnLine() error {
	if err := this.fsm.Event("online"); err != nil {
		logger.Errorf("farmer online return err: %v", err)
		return err
	}

	return nil
}

func (this *FarmerAccountHandler) Ping() error {
	if err := this.fsm.Event("ping"); err != nil {
		logger.Errorf("farmer ping return err: %v", err)
		return err
	}

	return nil
}

func (this *FarmerAccountHandler) Lost() error {
	if this.lostCount <= 0 {
		return errors.New("current lost count <= 0")
	}

	if err := this.fsm.Event("lost"); err != nil {
		logger.Errorf("farmer lost return err: %v", err)
		return err
	}

	return nil
}

func (this *FarmerAccountHandler) OffLine() error {
	if err := this.fsm.Event("offline"); err != nil {
		logger.Errorf("farmer offline return err: %v", err)
		return err
	}

	return nil
}

func (this *FarmerAccountHandler) beforeEvent(e *fsm.Event) {
	logger.Debugf("farmer(%v) do event(%v) from state %v enter state %v", this.account.FarmerID, e.Event, e.Src, e.Dst)
}

// after online, we set farmer's lost count 0
func (this *FarmerAccountHandler) afterOnLine() {
	npInterval := nextPingInterval()
	this.lostCount = 0
	this.nextPingTime = time.Now().Add(npInterval).UnixNano()

	time.AfterFunc(npInterval, func() {
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
}

func (this *FarmerAccountHandler) afterEvent(e *fsm.Event) {
	this.account.LastModifiedTime = time.Now().UnixNano()
	this.account.State = pb.FarmerState(pb.FarmerState_value[e.Dst])
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
	up := viper.GetInt("farmer.ping.up")
	down := viper.GetInt("farmer.ping.down")

	interval := basicInterval
	upflag := (rand.Int() % 2) == 0
	if upflag {
		interval += rand.Intn(up)
	} else {
		interval -= rand.Intn(down)
	}

	if interval < 0 {
		interval = basicInterval
	}

	return time.Duration(interval) * time.Second
}
