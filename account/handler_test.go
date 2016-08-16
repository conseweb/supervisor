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
	"gopkg.in/check.v1"
)

func (this *TestFarmerAccount) TestOnLine(c *check.C) {
	handler, _ := NewFarmerHandler("TestOnLine")
	c.Check(handler.OnLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_ONLINE)
}

func (this *TestFarmerAccount) TestLost(c *check.C) {
	handler, _ := NewFarmerHandler("TestLost")
	c.Check(handler.OnLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_ONLINE)

	handler.lostCount++
	c.Check(handler.Lost(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_LOST)
}

func (this *TestFarmerAccount) TestOffLine(c *check.C) {
	handler, _ := NewFarmerHandler("TestOffLine")
	c.Check(handler.OnLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_ONLINE)

	handler.lostCount++
	c.Check(handler.Lost(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_LOST)

	c.Check(handler.OffLine(), check.IsNil)
	c.Assert(handler.Account().State, check.Equals, pb.FarmerState_OFFLINE)
}
