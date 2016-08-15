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
	"github.com/conseweb/supervisor/api"
	pb "github.com/conseweb/supervisor/protos"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	default_addr = ":9376"
)

var (
	logger = logging.MustGetLogger("supervisor/node")
	server *grpc.Server
)

func StartNode() {
	defer func() {
		if x := recover(); x != nil {
			logger.Errorf("caught panic: %v", x)
		}
	}()

	addr := viper.GetString("node.address")
	if addr == "" {
		addr = default_addr
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("set up tcp listener err: %v", err)
	}

	grpc.EnableTracing = viper.GetBool("node.trace")
	server = grpc.NewServer()
	pb.RegisterFarmerPublicServer(server, &api.FarmerPublic{})

	logger.Infof("supervisor node listening on %s, waiting for connect...", addr)
	go server.Serve(lis)

	HandleNodeSignal()
}

// stop node
func StopNode() {
	server.GracefulStop()
}

func HandleNodeSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	logger.Infof("supervisor node has registered signal notify")

	for {
		s := <-sigs
		logger.Infof("supervisor node has received signal: %v", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info("supervisor node is graceful shutting down...")

			StopNode()

			logger.Info("supervisor node has exited")
			os.Exit(0)
		}
	}
}
