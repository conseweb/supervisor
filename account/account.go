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
	"sync"
	"time"

	"github.com/conseweb/supervisor/account/store"
	pb "github.com/conseweb/common/protos"
	"github.com/looplab/fsm"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"github.com/conseweb/supervisor/utils/semaphore"
	"github.com/conseweb/supervisor/challenge"
)

const (
	default_storage_backend = "rocksdb"
)

var (
	logger     *logging.Logger
	once       *sync.Once
	controller *FarmerAccountController
)

func init() {
	logger = logging.MustGetLogger("supervisor")
	once = &sync.Once{}
}

// read config file and init a storage backend
func getBackendStorage() (storage store.Storage) {
	backend := default_storage_backend
	if viper.GetString("account.store.backend") != "" {
		backend = viper.GetString("account.store.backend")
	}

	// use rocksdb as backend
	if backend == "rocksdb" {
		dbpath := viper.GetString("account.store.rocksdb.dbpath")
		if dbpath == "" {
			logger.Panic("using rocksdb as backend, but not specified dbpath")
		}

		var err error
		storage, err = store.NewStore(backend, dbpath)
		if err != nil {
			logger.Panic(err)
		}
	}

	return storage
}

type FarmerAccountController struct {
	accountStorage store.Storage
	accountTree    *AccountTree
	l              *sync.RWMutex
}

func getController() *FarmerAccountController {
	once.Do(func() {
		if controller != nil {
			return
		}

		controller = &FarmerAccountController{
			accountStorage: getBackendStorage(),
			accountTree:    NewAccountTree(),
			l:              &sync.RWMutex{},
		}

		go controller.checkHandlers()
	})

	return controller
}

// NewFarmer doesn't mean the farmer is already online
// just stands for there is a farmer want to connect 2 supervisor
func NewFarmerHandler(farmerId string) (*FarmerAccountHandler, error) {
	return getController().NewFarmerHandler(farmerId)
}
func (ctr *FarmerAccountController) NewFarmerHandler(farmerId string) (handler *FarmerAccountHandler, err error) {
	if farmerId == "" {
		err = errors.New("farmerId is empty")
		return
	}

	key := farmerId2Key(farmerId)
	// 1. looking farmer from account tree
	{
		handler = &FarmerAccountHandler{}
		ctr.l.RLock()
		handler, err = ctr.accountTree.Get(key)
		ctr.l.RUnlock()
		if err == nil {
			return
		}
	}

	tmpFsm := fsm.NewFSM(pb.FarmerState_OFFLINE.String(), fsm.Events{
		{Name: "offline", Src: []string{pb.FarmerState_ONLINE.String(), pb.FarmerState_LOST.String()}, Dst: pb.FarmerState_OFFLINE.String()},
		{Name: "online", Src: []string{pb.FarmerState_OFFLINE.String(), pb.FarmerState_LOST.String()}, Dst: pb.FarmerState_ONLINE.String()},
		{Name: "lost", Src: []string{pb.FarmerState_ONLINE.String()}, Dst: pb.FarmerState_LOST.String()},
	}, fsm.Callbacks{
		"before_event": func(e *fsm.Event) {
			handler.beforeEvent(e)
		},
	})

	// 2. looking farmer from storage, if found, put into account tree, and create fsm
	{
		handler = &FarmerAccountHandler{}
		var farmerBytes []byte
		ctr.l.RLock()
		farmerBytes, err = ctr.accountStorage.Get([]byte(key))
		ctr.l.RUnlock()
		if err == nil {
			handler.account, err = bytes2FarmerAccount(farmerBytes)
			if err == nil {
				handler.fsm = tmpFsm
				handler.account.State = pb.FarmerState(pb.FarmerState_value[handler.fsm.Current()])

				ctr.l.Lock()
				// put into account tree
				ctr.accountTree.Put(key, handler)
				ctr.l.Unlock()

				return
			}
		}
	}

	// 3 if can not load farmer account info from tree & storage, new a farmer account info
	{
		handler = &FarmerAccountHandler{}
		ctr.l.Lock()
		handler.account = &pb.FarmerAccount{
			FarmerID:         farmerId,
			Balance:          0,
			State:            pb.FarmerState_OFFLINE,
			LastModifiedTime: time.Now().UnixNano(),
		}
		handler.fsm = tmpFsm
		// put into account tree
		ctr.accountTree.Put(key, handler)
		if farmerBytes, err := farmerAccount2Bytes(handler.account); err == nil {
			ctr.asyncPersistFarmerBytes([]byte(key), farmerBytes)
		}
		ctr.l.Unlock()

		return
	}
}

func UpdateFarmerHandler(handler *FarmerAccountHandler) {
	getController().UpdateFarmerHandler(handler)
}

func (ctr *FarmerAccountController) UpdateFarmerHandler(handler *FarmerAccountHandler) {
	key := farmerId2Key(handler.account.FarmerID)

	ctr.l.Lock()
	if farmerBytes, err := farmerAccount2Bytes(handler.account); err == nil {
		// save back 2 memory
		if handler.account.State != pb.FarmerState_OFFLINE {
			ctr.accountTree.Put(key, handler)
		} else {
			ctr.accountTree.Delete(key)
		}

		ctr.asyncPersistFarmerBytes([]byte(key), farmerBytes)
	}
	ctr.l.Unlock()
}

// close the backend storage
func Close() error {
	return getController().Close()
}
func (ctr *FarmerAccountController) Close() error {
	ctr.accountTree = nil
	return ctr.accountStorage.Close()
}

// save back 2 storage, async
func (ctr *FarmerAccountController) asyncPersistFarmerBytes(farmerKey, farmerBytes []byte) {
	go ctr.accountStorage.Set(farmerKey, farmerBytes)
}

// check handlers
func (ctr *FarmerAccountController) checkHandlers() {
	ticker := time.NewTicker(getControllerCheckInterval())
	sema := semaphore.NewSemaphore(getControllerCheckWorkers())

	for {
		select {
		case <-ticker.C:
			logger.Infof("current farmer handler count: %d", ctr.accountTree.Len())

			for _, key := range ctr.accountTree.Keys() {
				sema.Acquire()
				go func(key string) {
					defer sema.Release()

					h, err := ctr.accountTree.Get(key)
					if err != nil {
						logger.Errorf("get farmer handler err: %v", err)
						return
					}

					// if handler's nextPingTime is before now, lostcount ++
					if time.Unix(h.nextPingTime, 0).Before(time.Now()) {
						h.lostCount++
						h.Lost()

						if h.lostCount >= viper.GetInt("farmer.ping.lostcount") {
							h.OffLine()
						}

						ctr.UpdateFarmerHandler(h)
					}

					// if handler's nextConquerTime > 0, nextChallengeReq isn't nil, and is before now, punlish
					if h.nextConquerTime > 0 && h.nextFarmerChallengeReq != nil && time.Unix(h.nextConquerTime, 0).Before(time.Now()) {
						blocksRange := h.nextFarmerChallengeReq.BlocksRange()
						challenge.GetFarmerChallengeReqCache().DelFarmerChallengeReq(h.nextFarmerChallengeReq.FarmerID(), blocksRange.HighBlockNumber, blocksRange.LowBlockNumber, h.nextFarmerChallengeReq.HashAlgo())

						h.nextConquerTime = 0
						h.nextFarmerChallengeReq = nil
						h.punishBalance()

						ctr.UpdateFarmerHandler(h)
					}
				}(key)
			}
		}
	}
}

func getControllerCheckInterval() time.Duration {
	if interval, err := time.ParseDuration(viper.GetString("account.check.interval")); err ==nil {
		return interval
	}

	viper.Set("account.check.interval", "60s")
	return time.Duration(60) * time.Second
}

func getControllerCheckWorkers() int {
	workers := viper.GetInt("account.check.workers")
	if workers <= 0 {
		viper.Set("account.check.workers", 8)
		workers = 8
	}

	return workers
}