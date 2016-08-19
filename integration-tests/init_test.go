package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
	"gopkg.in/check.v1"
)

const (
	defaultTimeout = time.Second * 3
	farmerID       = "hello"
)

var (
	SupervisorAddr string
)

type FarmerSuite struct {
	conn *grpc.ClientConn
}

func Test(t *testing.T) {
	check.TestingT(t)
}

func init() {
	SupervisorAddr = os.Getenv("SUPERVISOR_ADDR")
	if SupervisorAddr == "" {
		fmt.Println("SUPERVISOR_ADDR required.")
		os.Exit(1)
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithTimeout(defaultTimeout),
		grpc.WithBlock(),
	}

	conn, err := grpc.Dial(SupervisorAddr, opts...)
	if err != nil {
		fmt.Printf("dial grpc server err: %v", err)
		os.Exit(2)
	}

	fs := &FarmerSuite{conn}

	check.Suite(fs)
}
