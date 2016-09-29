//    Copyright 2016 Alex Macleod
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	yaml "gopkg.in/yaml.v2"
	"macleod.io/bounce/networking/client"
	"macleod.io/bounce/networking/network"
)

type Config struct {
	Version  int
	Servers  []*client.Server
	Networks []*network.Network
}

var (
	location = flag.String("config", getDefaultConfig(), "read configuration from a specified file")
	// Current is the current configuration, may not be saved to disk
	Current Config
)

// Load the configuration file into Current
func Load() {
	readConfigFile()

	tb, _ := yaml.Marshal(&Current)
	fmt.Println(string(tb))
}

func readConfigFile() {
	bytes, err := ioutil.ReadFile(*location)
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(bytes, &Current)
}

func getDefaultConfig() string {
	if runtime.GOOS == "windows" {
		return os.ExpandEnv("$USERPROFILE\\.bounce\\test.yaml")
	}
	if os.ExpandEnv("HOME") != "" {
		return os.ExpandEnv("$HOME/.bounce/test.yaml")
	}
	log.Fatal("Couldn't find home directory")
	return ""
}
