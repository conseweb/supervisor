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
	"github.com/conseweb/supervisor/account/store"
	"github.com/conseweb/supervisor/account/tree"
	pb "github.com/conseweb/supervisor/protos"
	"github.com/looplab/fsm"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"sync"
	"time"
	"github.com/golang/protobuf/proto"
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
	logger = logging.MustGetLogger("supervisor/account")
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

func getController() *FarmerAccountController {
	once.Do(func() {
		if controller != nil {
			return
		}

		controller = &FarmerAccountController{
			accountStorage: getBackendStorage(),
			accountTree:    tree.NewTree(),
			farmerFSMs:     make(map[string]*fsm.FSM),
			l:              &sync.RWMutex{},
		}
	})

	return controller
}

func NewFarmerHandler(farmerId string) (handler *FarmerAccountHandler) {
	return getController().NewFarmerHandler(farmerId)
}

func Close() error {
	return getController().Close()
}

type FarmerAccountController struct {
	accountStorage store.Storage
	accountTree    tree.Tree
	farmerFSMs     map[string]*fsm.FSM
	l              *sync.RWMutex
}

// NewFarmer doesn't mean the farmer is already online
// just stands for there is a farmer want to connect 2 supervisor
func (this *FarmerAccountController) NewFarmerHandler(farmerId string) (handler *FarmerAccountHandler) {
	key := farmerId2Key(farmerId)

	// 1. looking farmer from account tree
	{
		this.l.RLock()
		farmerBytes, err := this.accountTree.Get(key)
		this.l.RUnlock()
		if err == nil {
			// means farmer already in memory(tree), setup a farmer account handler
			this.l.RLock()
			fsm, ok := this.farmerFSMs[key]
			this.l.RUnlock()

			if ok {
				handler = newFarmerAccountHandler(bytes2FarmerAccount(farmerBytes), fsm)
				return
			}
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
		"after_online": func(e *fsm.Event) {
			handler.afterOnLine()
		},
		"after_event": func(e *fsm.Event) {
			handler.afterEvent(e)

			farmerBytes := farmerAccount2Bytes(handler.account)
			// save back 2 memory
			this.accountTree.Put(key, farmerBytes)
			// save back 2 storage, async
			go this.accountStorage.Set([]byte(key), farmerBytes)
		},
	})

	// 2. looking farmer from storage, if found, put into account tree, and create fsm
	{
		this.l.RLock()
		farmerBytes, err := this.accountStorage.Get([]byte(key))
		this.l.RUnlock()
		if err == nil {
			this.l.Lock()
			handler = newFarmerAccountHandler(bytes2FarmerAccount(farmerBytes), tmpFsm)
			// put into account tree
			this.accountTree.Put(key, farmerBytes)
			// put handler'fsm into fsm's map
			this.farmerFSMs[key] = handler.fsm
			this.l.Unlock()

			return
		}
	}

	// 3 if can not load farmer account info from tree & storage, new a farmer account info
	{
		this.l.Lock()
		farmerAccount := &pb.FarmerAccount{
			FarmerID:         farmerId,
			Balance:          0,
			State:            pb.FarmerState_OFFLINE,
			LastModifiedTime: time.Now().UnixNano(),
		}

		farmerBytes := farmerAccount2Bytes(farmerAccount)
		handler = newFarmerAccountHandler(farmerAccount, tmpFsm)
		// put into account tree
		this.accountTree.Put(key, farmerBytes)
		// put into account storage, async
		go this.accountStorage.Set([]byte(key), farmerBytes)

		// put handler'fsm into fsm's map
		this.farmerFSMs[key] = handler.fsm
		this.l.Unlock()

		return
	}
}

// close the backend storage
//
func (this *FarmerAccountController) Close() error {
	this.farmerFSMs = nil
	return this.accountStorage.Close()
}

type FarmerAccountHandler struct {
	account   *pb.FarmerAccount
	fsm       *fsm.FSM
	lostCount uint8
}

func (this *FarmerAccountHandler) OnLine() (*pb.FarmerAccount, error) {
	if err := this.fsm.Event("online"); err != nil {
		logger.Errorf("farmer online return err: %v", err)
		return nil, err
	}

	logger.Debug("alreay online")
	return this.account, nil
}

func (this *FarmerAccountHandler) Ping() {
	
}

func (this *FarmerAccountHandler) beforeEvent(e *fsm.Event) {
	logger.Debugf("Farmer[%v] enter state %v from state %v", this.account.FarmerID, e.Dst, e.Src)
}

// after online, we set farmer's lost count 0
func (this *FarmerAccountHandler) afterOnLine() {
	this.lostCount = 0
	logger.Debugf("i am in afterOnLine")
}

func (this *FarmerAccountHandler) afterEvent(e *fsm.Event) {
	this.account.LastModifiedTime = time.Now().UnixNano()
	this.account.State = pb.FarmerState(pb.FarmerState_value[e.Dst])
	logger.Debugf("i am in afterEvent")
}

func newFarmerAccountHandler(account *pb.FarmerAccount, fsm *fsm.FSM) *FarmerAccountHandler {
	if account.FarmerID == "" {
		return nil
	}

	return &FarmerAccountHandler{
		account: account,
		fsm:     fsm,
	}
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
func farmerAccount2Bytes(farmerObj *pb.FarmerAccount) []byte {
	bytesInfo, err := proto.Marshal(farmerObj)
	if err != nil {
		logger.Errorf("marshal farmer account err: %v", err)
		return nil
	}

	return bytesInfo
}

// unmarshal proto bytes 2 farmer account
// if there is an error, return nil
func bytes2FarmerAccount(fBytes []byte) *pb.FarmerAccount {
	farmerObj := &pb.FarmerAccount{}
	if err := proto.Unmarshal(fBytes, farmerObj); err != nil {
		logger.Errorf("unmarshal farmer account err: %v", err)
		return nil
	}

	return farmerObj
}
