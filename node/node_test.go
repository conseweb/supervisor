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
package node

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/conseweb/supervisor/challenge"
	"github.com/conseweb/supervisor/cli"
	pb "github.com/conseweb/supervisor/protos"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type NodeTest struct {
	conn   *grpc.ClientConn
	client pb.FarmerPublicClient
}

var _ = check.Suite(&NodeTest{})

func (t *NodeTest) SetUpSuite(c *check.C) {
	loadConfigure()

	go StartNode()
	time.Sleep(time.Second)
}

func (t *NodeTest) SetUpTest(c *check.C) {
	conn, err := cli.NewClientConnectionWithAddress(viper.GetString("node.address"), true, false, nil)
	c.Check(err, check.IsNil)

	t.conn = conn
	t.client = pb.NewFarmerPublicClient(t.conn)
}

func (t *NodeTest) TearDownTest(c *check.C) {
	t.conn.Close()
	t.client = nil
}

func (t *NodeTest) TearDownSuite(c *check.C) {
	StopNode()
}

func (t *NodeTest) TestFarmerOnLine(c *check.C) {
	farmerId := "TestFarmerOnLine"

	rsp, err := t.client.FarmerOnLine(context.Background(), &pb.FarmerOnLineReq{
		FarmerID: farmerId,
	})
	c.Check(err, check.IsNil)
	c.Check(rsp, check.NotNil)
	c.Check(rsp.GetError().OK(), check.Equals, true)
}

func (t *NodeTest) TestFarmerPing(c *check.C) {
	farmerId := "TestFarmerPing"

	rspPing, err := t.client.FarmerPing(context.Background(), &pb.FarmerPingReq{
		FarmerID: farmerId,
		BlocksRange: &pb.BlocksRange{
			HighBlockNumber: 100,
			LowBlockNumber:  10,
		},
	})
	c.Check(err, check.IsNil)
	c.Check(rspPing.GetError().OK(), check.Equals, false)

	rspOnLine, err := t.client.FarmerOnLine(context.Background(), &pb.FarmerOnLineReq{
		FarmerID: farmerId,
	})
	c.Check(err, check.IsNil)
	c.Check(rspOnLine, check.NotNil)
	c.Check(rspOnLine.GetError().OK(), check.Equals, true)

	// ten times ping
	for i := 0; i < 10; i++ {
		rspPing, err := t.client.FarmerPing(context.Background(), &pb.FarmerPingReq{
			FarmerID: farmerId,
			BlocksRange: &pb.BlocksRange{
				HighBlockNumber: 100,
				LowBlockNumber:  10,
			},
		})
		c.Check(err, check.IsNil)
		c.Check(rspPing.GetError().OK(), check.Equals, true)
	}
}

func (t *NodeTest) TestFarmerConquerChallenge(c *check.C) {
	farmerId := "TestFarmerConquerChallenge"

	rspOnLine, err := t.client.FarmerOnLine(context.Background(), &pb.FarmerOnLineReq{
		FarmerID: farmerId,
	})
	c.Check(err, check.IsNil)
	c.Check(rspOnLine, check.NotNil)
	c.Check(rspOnLine.GetError().OK(), check.Equals, true)

	// ten times ping
	for i := 0; i < 10; i++ {
		rspPing, err := t.client.FarmerPing(context.Background(), &pb.FarmerPingReq{
			FarmerID: farmerId,
			BlocksRange: &pb.BlocksRange{
				HighBlockNumber: 100,
				LowBlockNumber:  10,
			},
		})
		c.Check(err, check.IsNil)
		c.Check(rspPing.GetError().OK(), check.Equals, true)

		if rspPing.NeedChallenge {
			rspChallenge, err := t.client.FarmerConquerChallenge(context.Background(), &pb.FarmerConquerChallengeReq{
				FarmerID:    farmerId,
				BlocksHash:  challenge.FarmerBindConquerHash(farmerId, rspPing.HashAlgo, challenge.HASH(rspPing.HashAlgo, []byte(""))),
				HashAlgo:    rspPing.HashAlgo,
				BlocksRange: rspPing.GetBlocksRange(),
			})

			c.Check(err, check.IsNil)
			c.Check(rspChallenge.GetError().OK(), check.Equals, true)
		}
	}
}

func (t *NodeTest) TestFarmerOffLine(c *check.C) {
	farmerId := "TestFarmerOffLine"

	rspOnLine, err := t.client.FarmerOnLine(context.Background(), &pb.FarmerOnLineReq{
		FarmerID: farmerId,
	})
	c.Check(err, check.IsNil)
	c.Check(rspOnLine, check.NotNil)
	c.Check(rspOnLine.GetError().OK(), check.Equals, true)

	rspOffLine, err := t.client.FarmerOffLine(context.Background(), &pb.FarmerOffLineReq{
		FarmerID: farmerId,
	})

	c.Check(err, check.IsNil)
	c.Check(rspOffLine, check.NotNil)
	c.Check(rspOffLine.GetError().OK(), check.Equals, true)
	c.Check(rspOffLine.Account.State, check.Equals, pb.FarmerState_OFFLINE)
}

func loadConfigure() {
	// Now set the configuration file
	viper.SetEnvPrefix(strings.ToUpper("tcsv"))
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetConfigName("config") // name of config file (without extension)

	alternativeCfgPath := os.Getenv("TCSV_CFG_PATH")
	if alternativeCfgPath != "" {
		logger.Info("User defined config file path: %s", alternativeCfgPath)
		viper.AddConfigPath(alternativeCfgPath)
	} else {
		viper.AddConfigPath("./")
		for _, p := range filepath.SplitList(os.Getenv("GOPATH")) {
			viper.AddConfigPath(filepath.Join(p, "src/github.com/conseweb/supervisor"))
		}
	}

	viper.AddConfigPath("./")   // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logger.Panicf("Fatal error config file: %s \n", err)
	}
}
