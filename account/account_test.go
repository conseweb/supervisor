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
	pb "github.com/conseweb/supervisor/protos"
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
}

func (this *TestFarmerAccount) TearDownSuite(c *check.C) {
	time.Sleep(time.Second)
	Close()
	os.Remove(filepath.Join(os.TempDir(), "testAccount"))
}

func (this *TestFarmerAccount) TestOnLine(c *check.C) {
	handler := NewFarmerHandler("farmerId0001")
	account, err := handler.OnLine()
	c.Check(err, check.IsNil)
	c.Check(account, check.NotNil)
	c.Assert(account.FarmerID, check.DeepEquals, "farmerId0001")
}

func (this *TestFarmerAccount) TestFarmerMarshalUnmarshal(c *check.C) {
	farmer := &pb.FarmerAccount{
		FarmerID: "1234567",
		Balance:  100,
	}

	farmerBytes := farmerAccount2Bytes(farmer)
	tmpFarmer := bytes2FarmerAccount(farmerBytes)

	c.Assert(tmpFarmer.FarmerID, check.Equals, farmer.FarmerID)
	c.Assert(tmpFarmer.Balance, check.Equals, farmer.Balance)
}

func (this *TestFarmerAccount) BenchmarkFarmerMarshal(c *check.C) {
	for i := 0; i < c.N; i++ {
		farmerAccount2Bytes(&pb.FarmerAccount{
			FarmerID: "1234567",
			Balance:  100,
		})
	}
}

func (this *TestFarmerAccount) BenchmarkFarmerUnmarshal(c *check.C) {
	fBytes := farmerAccount2Bytes(&pb.FarmerAccount{
		FarmerID: "1234567",
		Balance:  100,
	})
	for i := 0; i < c.N; i++ {
		bytes2FarmerAccount(fBytes)
	}
}
