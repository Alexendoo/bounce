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

package client

import (
	"bufio"
	"net"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func TestRegistration(t *testing.T) {
	network := &Network{
		Name: "Test",

		Host: "localhost",
		Port: "25000",

		Nick: "Nick",
		User: "User",
		Real: "Real name",
	}

	done := make(chan bool)

	go func() {
		listener, _ := net.Listen("tcp", "localhost:25000")
		conn, _ := listener.Accept()

		scanner := bufio.NewScanner(conn)
		c.Convey("Requests capabilities", t, func() {
			scanner.Scan()
			c.So(scanner.Text(), c.ShouldEqual, "CAP LS 302")
		})
		c.Convey("Sets nickname", t, func() {
			scanner.Scan()
			c.So(scanner.Text(), c.ShouldEqual, "NICK Nick")
		})
		c.Convey("Sets username", t, func() {
			scanner.Scan()
			c.So(scanner.Text(), c.ShouldEqual, "USER User - - :Real name")
		})
		done <- true
	}()

	network.connect()
	network.register()
	<-done
	network.disconnect()
}
