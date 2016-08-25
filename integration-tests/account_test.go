package main

import (
	pb "github.com/conseweb/common/protos"
	"golang.org/x/net/context"
	"gopkg.in/check.v1"
)

func (f *FarmerSuite) TestAccount(c *check.C) {
	c.Assert(f.conn, check.NotNil)
	client := pb.NewFarmerPublicClient(f.conn)

	onRes := onlineAccount(c, client)
	c.Assert(onRes.Error.ErrorType, check.Equals, pb.ErrorType(pb.ErrorType_INVALID_STATE_FARMER_ONLINE))

	onlineAccount(c, client)
	c.Assert(onRes.Error.ErrorType, check.Equals, pb.ErrorType(pb.ErrorType_INVALID_STATE_FARMER_ONLINE))
	defer offlineAccount(c, client)
	pingAccount(c, client)

}

func onlineAccount(c *check.C, client pb.FarmerPublicClient) *pb.FarmerOnLineRsp {
	onlineReq := &pb.FarmerOnLineReq{FarmerID: "hello"}
	res, err := client.FarmerOnLine(context.Background(), onlineReq)
	c.Assert(err, check.IsNil)
	c.Assert(res, check.NotNil)

	c.Logf("online error: %+v", res.Error)
	c.Logf("online account: %+v", res.Account)

	return res
}

func pingAccount(c *check.C, client pb.FarmerPublicClient) *pb.FarmerPingRsp {
	pingReq := &pb.FarmerPingReq{
		FarmerID: farmerID,
		BlocksRange: &pb.BlocksRange{
			LowBlockNumber:  uint64(1),
			HighBlockNumber: uint64(120),
		},
	}
	res, err := client.FarmerPing(context.Background(), pingReq)
	c.Assert(err, check.IsNil)
	c.Assert(res, check.NotNil)

	c.Logf("error: %+v", res.Error)
	c.Assert(res.Error.ErrorType, check.Equals, pb.ErrorType(pb.ErrorType_NONE_ERROR))
	return res
}

func offlineAccount(c *check.C, client pb.FarmerPublicClient) {
	offlineReq := &pb.FarmerOffLineReq{FarmerID: "hello"}
	res, err := client.FarmerOffLine(context.Background(), offlineReq)
	c.Assert(err, check.IsNil)
	c.Assert(res, check.NotNil)
}
