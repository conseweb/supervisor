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
package main

import (
	"github.com/conseweb/common/config"
	"github.com/conseweb/supervisor/node"
	"github.com/hyperledger/fabric/flogging"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

const (
	appName = "supervisor"
)

var (
	logger = logging.MustGetLogger("main")
	app    = kingpin.New(appName, "A command-line trust-chain supervisor cli.")
	svnode = app.Command("node", "Supervisor Node")
)

func init() {
	// load configure
	if err := config.LoadConfig("SUPERVISOR", "supervisor", "github.com/conseweb/supervisor"); err != nil {
		// Handle errors reading the config file
		logger.Panicf("Fatal error config file: %s \n", err)
	}

	flogging.LoggingInit("main")
}

func main() {
	app.Version(viper.GetString("server.version"))
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case svnode.FullCommand():
		node.StartNode()
	}
}
