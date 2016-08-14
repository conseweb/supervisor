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
	"golang.org/x/net/context"
)

type fmp struct {
}

func (this *fmp) FarmerOnLine(ctx context.Context, req *pb.FarmerOnLineReq) (*pb.FarmerOnLineRsp, error) {
	// TODO vairfy farmerid

	handler := account.NewFarmerHandler(req.FarmerID)

	// online
	if err := handler.OnLine(); err != nil {
		return nil, err
	}

	return &pb.FarmerOnLineRsp{
		Account:  handler.Account(),
		NextPing: handler.NextPingTime(),
	}, nil
}

func (this *fmp) FarmerPing(ctx context.Context, req *pb.FarmerPingReq) (*pb.FarmerPingRsp, error) {
	// TODO vairfy farmerid
	handler := account.NewFarmerHandler(req.FarmerID)
	rsp := &pb.FarmerPingRsp{}

	// online event
	if err := handler.OnLine(); err != nil {
		return nil, err
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

func (this *fmp) FarmerConquerChallenge(ctx context.Context, req *pb.FarmerConquerChallengeReq) (*pb.FarmerConquerChallengeRsp, error) {
	handler := account.NewFarmerHandler(req.FarmerID)
	rsp := &pb.FarmerConquerChallengeRsp{
		Account: handler.Account(),
	}

	if err := handler.ConquerChallenge(req.BlocksRange.HighBlockNumber, req.BlocksRange.LowBlockNumber, req.HashAlgo, req.BlocksHash); err != nil {
		rsp.ConquerOK = false
		return rsp, err
	}
	account.UpdateFarmerHandler(handler)

	rsp.ConquerOK = true
	return rsp, nil
}

func (this *fmp) FarmerOffLine(ctx context.Context, req *pb.FarmerOffLineReq) (*pb.FarmerOffLineRsp, error) {

	// TODO varify farmer id
	handler := account.NewFarmerHandler(req.FarmerID)

	// offline event
	if err := handler.OffLine(); err != nil {
		return nil, err
	}

	return &pb.FarmerOffLineRsp{
		Account: handler.Account(),
	}, nil
}
