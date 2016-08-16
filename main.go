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
package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/conseweb/supervisor/node"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	appName = "tcsv"
)

var (
	logger = logging.MustGetLogger("supervisor")

	app = kingpin.New(appName, "A command-line trust-chain supervisor cli.")

	svnode = app.Command("node", "Supervisor Node")
)

func init() {
	// load configure
	loadConfigure()

	level, err := logging.LogLevel(viper.GetString("server.logging"))
	if err != nil {
		logger.Fatalf("set logging level err: %v", level)
	}

	logging.SetLevel(level, "supervisor")
}

func main() {

	app.Version(viper.GetString("server.version"))
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case svnode.FullCommand():
		node.StartNode()
	}
}

func loadConfigure() {
	// Now set the configuration file
	viper.SetEnvPrefix(strings.ToUpper(appName))
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetConfigName("config") // name of config file (without extension)

	alternativeCfgPath := os.Getenv("TCSV_CFG_PATH")
	if alternativeCfgPath != "" {
		logger.Info("User defined config file path: %s", alternativeCfgPath)
		viper.AddConfigPath(alternativeCfgPath)
	} else {
		viper.AddConfigPath("./")
		for _, p := range filepath.SplitList(os.Getenv("GOPATH")) {
			viper.AddConfigPath(filepath.Join(p, "src/github.com/conseweb/supervisor"))
		}
	}

	viper.AddConfigPath("./")   // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logger.Panicf("Fatal error config file: %s \n", err)
	}
}
