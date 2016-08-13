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
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAccount(t *testing.T) {
	check.TestingT(t)
}

type TestFarmerAccount struct {
}

var _ = check.Suite(&TestFarmerAccount{})

func (this *TestFarmerAccount) SetUpSuite(c *check.C) {
	viper.Set("account.store.backend", "rocksdb")
	viper.Set("account.store.rocksdb.dbpath", filepath.Join(os.TempDir(), "testAccount"))
	viper.Set("farmer.ping.interval", 900)
	viper.Set("farmer.ping.up", 900)
	viper.Set("farmer.ping.down", 800)
	viper.Set("farmer.ping.lostcount", 2)
	viper.Set("farmer.challenge.hash", "SHA256")
}

func (this *TestFarmerAccount) TearDownSuite(c *check.C) {
	time.Sleep(time.Second)
	Close()
	os.Remove(filepath.Join(os.TempDir(), "testAccount"))
}
