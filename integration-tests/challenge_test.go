package main

import (
	"time"

	pb "github.com/conseweb/supervisor/protos"
	"golang.org/x/net/context"
	"gopkg.in/check.v1"
)

func (f *FarmerSuite) TestChallenge(c *check.C) {
	c.Assert(f.conn, check.NotNil)
	client := pb.NewFarmerPublicClient(f.conn)

	onlineAccount(c, client)
	defer offlineAccount(c, client)
	for i := 0; i < 20; i++ {
		res := pingAccount(c, client)
		if !res.NeedChallenge {
			time.Sleep(1 * time.Second)
			continue
		}
		return
	}
}

func challengeAccount(c *check.C, client pb.FarmerPublicClient, hstr string, hAlgo int32, br *pb.BlocksRange) {
	challengeReq := &pb.FarmerConquerChallengeReq{
		FarmerID:    farmerID,
		BlocksHash:  hstr,
		HashAlgo:    pb.HashAlgo(hAlgo),
		BlocksRange: br,
	}

	res, err := client.FarmerConquerChallenge(context.Background(), challengeReq)
	c.Assert(err, check.IsNil)
	c.Assert(res, check.NotNil)

	c.Logf("error: %+v", res.Error)
	c.Assert(res.Error.ErrorType, check.Equals, pb.ErrorType(0))

	c.Assert(res.ConquerOK, check.Equals, true)
}
