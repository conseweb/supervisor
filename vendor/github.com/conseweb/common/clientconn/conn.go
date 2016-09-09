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
package clientconn

import (
	"github.com/hyperledger/fabric/flogging"
	"github.com/op/go-logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"sync"
	"time"
)

const defaultTimeout = time.Second * 3

var (
	once   = &sync.Once{}
	logger = logging.MustGetLogger("cli")
)

// NewClientConnectionWithAddress Returns a new grpc.ClientConn to the given address.
func NewClientConnectionWithAddress(address string, block bool, tslEnabled bool, creds credentials.TransportAuthenticator) (*grpc.ClientConn, error) {
	once.Do(func() {
		flogging.LoggingInit("cli")
	})

	var opts []grpc.DialOption
	if tslEnabled {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithTimeout(defaultTimeout))
	if block {
		opts = append(opts, grpc.WithBlock())
	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		logger.Errorf("dial grpc server err: %v", err)
		return nil, err
	}
	return conn, nil
}

// InitTLSForClient returns TLS credentials for client
func InitTLSForClient(hostoverride, certFile string) credentials.TransportAuthenticator {
	once.Do(func() {
		flogging.LoggingInit("cli")
	})

	creds, err := credentials.NewClientTLSFromFile(certFile, hostoverride)
	if err != nil {
		logger.Errorf("Failed to create TLS credentials %v", err)
		creds = credentials.NewClientTLSFromCert(nil, hostoverride)
	}

	return creds
}
