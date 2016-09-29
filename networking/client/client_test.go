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

	"macleod.io/bounce/irc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		message *irc.Message
		head    net.Conn
		tail    net.Conn
		client  *Client
	)

	BeforeEach(func() {
		message = &irc.Message{
			Prefix:  "example.org",
			Command: "PING",
		}
		head, tail = net.Pipe()
		client = New(head)
	})

	AfterEach(func() {
		Expect(client.Close()).NotTo(HaveOccurred())
	})

	It("Accepts messages", func() {
		client.In <- message
		scanner := bufio.NewScanner(tail)
		Expect(scanner.Scan()).To(BeTrue())
		received := irc.ParseMessage(scanner.Text())
		Expect(received.Prefix).To(Equal(message.Prefix))
		Expect(received.Command).To(Equal(message.Command))
	})

	It("Emits messages", func() {
		message.Buffer().WriteTo(tail)
		received := <-client.Out
		Expect(received.Prefix).To(Equal(message.Prefix))
		Expect(received.Command).To(Equal(message.Command))
	})
})
