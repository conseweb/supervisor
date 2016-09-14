/*
Copyright Mojing Inc. 2016 All Rights Reserved.
Written by mint.zhao.chiu@gmail.com. github.com: https://www.github.com/mintzhao

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
	"net"

	"github.com/conseweb/common/clientconn"
	"github.com/conseweb/common/exec"
	pb "github.com/conseweb/common/protos"
	"github.com/conseweb/supervisor/account"
	"github.com/conseweb/supervisor/api"
	"github.com/conseweb/supervisor/challenge"
	"github.com/hyperledger/fabric/flogging"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	default_addr = ":9376"
)

var (
	logger = logging.MustGetLogger("node")
	server *grpc.Server
)

func StartNode() {
	flogging.LoggingInit("node")

	// verify supervisor ok or not
	verifySupervisor()

	addr := viper.GetString("node.address")
	if addr == "" {
		addr = default_addr
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("set up tcp listener err: %v", err)
	}

	grpc.EnableTracing = viper.GetBool("node.trace")
	logger.Infof("grpc.EnableTracing: %v", grpc.EnableTracing)

	opts := []grpc.ServerOption{}
	if viper.GetBool("node.tls.enabled") {
		opts = append(opts, grpc.Creds(initTLSForServer()))
	}
	server = grpc.NewServer(opts...)

	// register
	flogging.LoggingInit("api")
	pb.RegisterFarmerPublicServer(server, &api.FarmerPublic{})

	go server.Serve(lis)
	logger.Infof("supervisor node listening on %s, waiting for connect...", addr)

	exec.HandleSignal(StopNode)
}

// stop node
func StopNode() error {
	// server.GracefulStop()
	if server != nil {
		server.Stop()
	}
	account.Close()
	challenge.GetFarmerChallengeReqCache().Close()
	challenge.GetBlocksHashCache().Close()

	return nil
}

// InitTLSForServer returns TLS credentials for node
func initTLSForServer() credentials.TransportAuthenticator {
	creds, err := credentials.NewServerTLSFromFile(viper.GetString("node.tls.cert.file"), viper.GetString("node.tls.key.file"))
	if err != nil {
		logger.Errorf("Failed to create TLS credentials %v", err)
		creds = credentials.NewServerTLSFromCert(nil)
	}

	return creds
}

func verifySupervisor() {
	logger.Info("begin to verify supervisor via idprovider")

	var conn *grpc.ClientConn
	var err error
	address := viper.GetString("idprovider.port")
	tlsEnable := viper.GetBool("idprovider.tls.enabled")
	if tlsEnable {
		hostoverride := viper.GetString("idprovider.tls.serverhostoverride")
		certFile := viper.GetString("idprovider.tls.cert.file")
		conn, err = clientconn.NewClientConnectionWithAddress(address, true, true, clientconn.InitTLSForClient(hostoverride, certFile))
	} else {
		conn, err = clientconn.NewClientConnectionWithAddress(address, true, false, nil)
	}
	if err != nil {
		logger.Fatalf("connect with idprovider return error: %v", err)
	}
	defer conn.Close()

	idpaCli := pb.NewIDPAClient(conn)
	rsp, err := idpaCli.VerifyDevice(context.Background(), &pb.VerifyDeviceReq{
		UserID:      viper.GetString("node.svorg"),
		DeviceID:    viper.GetString("node.svid"),
		DeviceAlias: viper.GetString("node.svalias"),
		For:         pb.DeviceFor_SUPERVISOR,
	})
	if err != nil {
		logger.Fatal(err)
	}

	if !rsp.Error.OK() {
		logger.Fatalf("varify supervisor via idprovider return error: %v", rsp.Error)
	}

	logger.Info("finish verify supervisor via idprovider, OK")
}
