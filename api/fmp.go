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
package api

import (
	"github.com/conseweb/supervisor/account"
	"github.com/conseweb/supervisor/challenge"
	pb "github.com/conseweb/supervisor/protos"
	"github.com/op/go-logging"
	"golang.org/x/net/context"
)

// TODO vairfy farmerid

var (
	logger = logging.MustGetLogger("supervisor")
)

type FarmerPublic struct {
}

func (this *FarmerPublic) FarmerOnLine(ctx context.Context, req *pb.FarmerOnLineReq) (*pb.FarmerOnLineRsp, error) {
	logger.Debugf("new connect for FarmerOnLine, req: %+v", req)

	rsp := &pb.FarmerOnLineRsp{
		Error: pb.ResponseOK(),
	}
	handler := account.NewFarmerHandler(req.FarmerID)

	// online
	if err := handler.OnLine(); err != nil {
		rsp.Error = pb.NewErrorf(pb.ErrorType_FARMER_ONLINE, "online return err: %v", err)
		return rsp, nil
	}

	rsp.Account = handler.Account()
	rsp.NextPing = handler.NextPingTime()

	return rsp, nil
}

func (this *FarmerPublic) FarmerPing(ctx context.Context, req *pb.FarmerPingReq) (*pb.FarmerPingRsp, error) {
	logger.Debugf("new connect for FarmerPing, req: %+v", req)

	rsp := &pb.FarmerPingRsp{
		Error: pb.ResponseOK(),
	}
	handler := account.NewFarmerHandler(req.FarmerID)

	if err := handler.Ping(); err != nil {
		rsp.Error = pb.NewError(pb.ErrorType_FARMER_ONLINE, err.Error())
		return rsp, nil
	}

	// need challenge
	need, brange := handler.NeedChallengeBlocks(req.BlocksRange.HighBlockNumber, req.BlocksRange.LowBlockNumber)
	rsp.NeedChallenge = need
	rsp.BlocksRange = brange
	rsp.Account = handler.Account()
	rsp.NextPing = handler.NextPingTime()
	if need {
		hashAlgo := handler.ChallengeHashAlgo()
		rsp.HashAlgo = hashAlgo

		// sv cache challenge req
		challenge.GetFarmerChallengeReqCache().SetFarmerChallengeReq(req.FarmerID, brange.HighBlockNumber, brange.LowBlockNumber, hashAlgo)
	}

	return rsp, nil
}

func (this *FarmerPublic) FarmerConquerChallenge(ctx context.Context, req *pb.FarmerConquerChallengeReq) (*pb.FarmerConquerChallengeRsp, error) {
	logger.Debugf("new connect for FarmerConquerChallenge, req: %+v", req)

	rsp := &pb.FarmerConquerChallengeRsp{
		Error: pb.ResponseOK(),
	}
	handler := account.NewFarmerHandler(req.FarmerID)

	if err := handler.ConquerChallenge(req.BlocksRange.HighBlockNumber, req.BlocksRange.LowBlockNumber, req.HashAlgo, req.BlocksHash); err != nil {
		rsp.ConquerOK = false
		rsp.Error = pb.NewErrorf(pb.ErrorType_FARMER_CHALLENGE_FAIL, "challenge fail: %v", err)

		return rsp, nil
	}

	rsp.ConquerOK = true
	rsp.Account = handler.Account()

	return rsp, nil
}

func (this *FarmerPublic) FarmerOffLine(ctx context.Context, req *pb.FarmerOffLineReq) (*pb.FarmerOffLineRsp, error) {
	logger.Debugf("new connect for FarmerOffLine, req: %+v", req)

	rsp := &pb.FarmerOffLineRsp{
		Error: pb.ResponseOK(),
	}
	handler := account.NewFarmerHandler(req.FarmerID)

	// offline event
	if err := handler.OffLine(); err != nil {
		rsp.Error = pb.NewErrorf(pb.ErrorType_FARMER_OFFLINE, "offline err: %v", err)
		return rsp, nil
	}

	rsp.Account = handler.Account()

	return rsp, nil
}
