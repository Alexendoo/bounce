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

package network_test

import (
	"bufio"
	"io"

	"macleod.io/bounce/irc"
	. "macleod.io/bounce/networking/network"

	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Network", func() {
	var (
		network  *Network
		listener net.Listener
		conn     net.Conn
		scanner  *bufio.Scanner
	)

	BeforeEach(func() {
		listener, _ = net.Listen("tcp", "localhost:0")
		network = &Network{
			Addr: listener.Addr().String(),
			Nick: "nickname",
			Real: "real name",
			User: "username",
		}
		Expect(network.Connect()).NotTo(HaveOccurred())
		var err error
		conn, err = listener.Accept()
		Expect(err).NotTo(HaveOccurred())
		scanner = bufio.NewScanner(conn)
	})

	AfterEach(func() {
		Expect(network.Close()).NotTo(HaveOccurred())
	})

	It("Should send initial registration messages", func() {
		go network.Register()

		Expect(scanner.Scan()).To(BeTrue())
		Expect(scanner.Text()).To(Equal("CAP LS 302"))
		Expect(scanner.Scan()).To(BeTrue())
		Expect(scanner.Text()).To(Equal("NICK nickname"))
		Expect(scanner.Scan()).To(BeTrue())
		Expect(scanner.Text()).To(Equal("USER username - - :real name"))
	})

	It("Should send incoming messages to the network", func() {
		network.In <- &irc.Message{
			Command: "PING",
		}
		Expect(scanner.Scan()).To(BeTrue())
		Expect(scanner.Text()).To(Equal("PING"))
	})

	It("Should emit messages from the network", func() {
		text := ":example.org PING\r\n"
		io.WriteString(conn, text)
		message := <-network.Out
		Expect(message.Buffer().String()).To(Equal(text))
	})

	It("Should return an error if it fails to connect", func() {
		network := &Network{
			Addr: "localhost:70000",
		}
		Expect(network.Connect()).To(HaveOccurred())
	})
})
