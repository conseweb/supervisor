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
package challenge

import (
	"fmt"
	pb "github.com/conseweb/supervisor/protos"
	"github.com/op/go-logging"
	"gopkg.in/check.v1"
)

type TestFarmerChallengeCache struct {
}

var _ = check.Suite(&TestFarmerChallengeCache{})

func (this *TestFarmerChallengeCache) SetUpSuite(c *check.C) {
	logging.SetLevel(logging.INFO, "supervisor/challenge")
}

func (this *TestFarmerChallengeCache) TestSetFarmerChallengeReq(c *check.C) {
	c.Check(GetFarmerChallengeReqCache().SetFarmerChallengeReq("farmerId001", 100, 20, pb.HashAlgo_SHA1), check.Equals, true)
	c.Check(GetFarmerChallengeReqCache().SetFarmerChallengeReq("farmerId001", 100, 20, pb.HashAlgo_SHA1), check.Equals, false)
}

func (this *TestFarmerChallengeCache) TestGetFarmerChallengeReq(c *check.C) {
	c.Check(GetFarmerChallengeReqCache().SetFarmerChallengeReq("farmerId002", 100, 20, pb.HashAlgo_SHA1), check.Equals, true)

	req, get := GetFarmerChallengeReqCache().GetFarmerChallengeReq("farmerId002", 100, 20, pb.HashAlgo_SHA1)
	c.Check(get, check.Equals, true)
	c.Check(req.farmerId, check.Equals, "farmerId002")
}

func (this *TestFarmerChallengeCache) TestDelFarmerChallengeReq(c *check.C) {
	c.Check(GetFarmerChallengeReqCache().SetFarmerChallengeReq("farmerId003", 100, 20, pb.HashAlgo_SHA1), check.Equals, true)

	req, get := GetFarmerChallengeReqCache().GetFarmerChallengeReq("farmerId003", 100, 20, pb.HashAlgo_SHA1)
	c.Check(get, check.Equals, true)
	c.Check(req.farmerId, check.Equals, "farmerId003")

	GetFarmerChallengeReqCache().DelFarmerChallengeReq("farmerId003", 100, 20, pb.HashAlgo_SHA1)

	req, get = GetFarmerChallengeReqCache().GetFarmerChallengeReq("farmerId003", 100, 20, pb.HashAlgo_SHA1)
	c.Check(get, check.Equals, false)
	c.Check(req, check.IsNil)
}

func (this *TestFarmerChallengeCache) BenchmarkSetFarmerChallengeReq(c *check.C) {
	for i := 0; i < c.N; i++ {
		GetFarmerChallengeReqCache().SetFarmerChallengeReq(fmt.Sprintf("farmerId%v", i), 100, 20, pb.HashAlgo_SHA1)
	}
}

func (this *TestFarmerChallengeCache) BenchmarkGetFarmerChallengeReq(c *check.C) {
	farmerId := "farmerIdGet"
	GetFarmerChallengeReqCache().SetFarmerChallengeReq(farmerId, 100, 20, pb.HashAlgo_SHA1)
	for i := 0; i < c.N; i++ {
		GetFarmerChallengeReqCache().GetFarmerChallengeReq(farmerId, 100, 20, pb.HashAlgo_SHA1)
	}
}

func (this *TestFarmerChallengeCache) BenchmarkDelFarmerChallengeReq(c *check.C) {
	farmerId := "farmerIdDel"
	GetFarmerChallengeReqCache().SetFarmerChallengeReq(farmerId, 100, 20, pb.HashAlgo_SHA1)
	for i := 0; i < c.N; i++ {
		GetFarmerChallengeReqCache().DelFarmerChallengeReq(farmerId, 100, 20, pb.HashAlgo_SHA1)
	}
}

