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

	"gopkg.in/yaml.v2"
)

type config struct {
	Version  int
	Name     string
	Networks []network
}

type network struct {
	Name string

	Address string
	Port    uint16

	Nick     string
	Realname string
	Username string
}

type password struct {
}

var (
	location string
	// Config is the current configuration, may not be saved to disk
	Config config
)

// Load the configuration file into Config
func Load() {
	parseFlags()
	readConfigFile()

	tb, _ := yaml.Marshal(&Config)
	fmt.Println(string(tb))
}

func parseFlags() {
	flag.StringVar(&location, "config", getConfigLocation(), "read configuration from a specified file")
	flag.Parse()
}

func readConfigFile() {
	location = getConfigLocation()
	bytes, err := ioutil.ReadFile(location)
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(bytes, &Config)
}

func getConfigLocation() string {
	if location != "" {
		return location
	}
	println(os.Getenv("HOME"))
	return "test.yaml"
}
