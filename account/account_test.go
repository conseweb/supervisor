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
	"os"
	"path/filepath"
	"testing"
	"time"

	pb "github.com/conseweb/supervisor/protos"
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
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
}

func (this *TestFarmerAccount) TearDownSuite(c *check.C) {
	time.Sleep(time.Second)
	Close()
	os.Remove(filepath.Join(os.TempDir(), "testAccount"))
}

func (this *TestFarmerAccount) TestOnLine(c *check.C) {
	handler := NewFarmerHandler("farmerId0001")
	c.Check(handler.OnLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_ONLINE)
}

func (this *TestFarmerAccount) TestLost(c *check.C) {
	handler := NewFarmerHandler("farmerId0002")
	c.Check(handler.OnLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_ONLINE)

	handler.lostCount++
	c.Check(handler.Lost(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_LOST)
}

func (this *TestFarmerAccount) TestOffLine(c *check.C) {
	handler := NewFarmerHandler("farmerId0003")
	c.Check(handler.OnLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_ONLINE)

	handler.lostCount++
	c.Check(handler.Lost(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_LOST)

	c.Check(handler.OffLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_OFFLINE)
}
