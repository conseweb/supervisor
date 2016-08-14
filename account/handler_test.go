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
	"github.com/conseweb/supervisor/challenge"
)

func (this *TestFarmerAccount) TestOnLine(c *check.C) {
	handler := NewFarmerHandler("TestOnLine")
	c.Check(handler.OnLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_ONLINE)
}

func (this *TestFarmerAccount) TestLost(c *check.C) {
	handler := NewFarmerHandler("TestLost")
	c.Check(handler.OnLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_ONLINE)

	handler.lostCount++
	c.Check(handler.Lost(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_LOST)
}

func (this *TestFarmerAccount) TestOffLine(c *check.C) {
	handler := NewFarmerHandler("TestOffLine")
	c.Check(handler.OnLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_ONLINE)

	handler.lostCount++
	c.Check(handler.Lost(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_LOST)

	c.Check(handler.OffLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_OFFLINE)
}

func (this *TestFarmerAccount) TestChallengeHashAlgo(c *check.C) {
	handler := NewFarmerHandler("TestChallengeHashAlgo")

	c.Assert(handler.ChallengeHashAlgo(), check.Equals, pb.HashAlgo_SHA256)

	viper.Set("farmer.challenge.hashalgo", "SHA512")
	c.Assert(handler.ChallengeHashAlgo(), check.Not(check.Equals), pb.HashAlgo_SHA256)
	c.Assert(handler.ChallengeHashAlgo(), check.Equals, pb.HashAlgo_SHA512)
}

func (this *TestFarmerAccount) TestConquerChallenge(c *check.C) {
	c.Skip("no blocks")
	handler := NewFarmerHandler("TestConquerChallenge")

	challenge.GetFarmerChallengeReqCache().SetFarmerChallengeReq("TestConquerChallenge", 100, 20, pb.HashAlgo_SHA256)
	c.Check(handler.ConquerChallenge(100, 20, pb.HashAlgo_SHA256, challenge.FarmerConquerHash("TestConquerChallenge", pb.HashAlgo_SHA256, challenge.HASH(pb.HashAlgo_SHA256, []byte("")))), check.IsNil)
}