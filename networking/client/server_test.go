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
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var server *Server
	var done chan bool

	BeforeEach(func() {
		server = &Server{
			Addr: "localhost:0",
		}
		done = make(chan bool)
	})

	AfterEach(func() {
		server.Close()
	})

	It("Accepts connections", func() {
		conns, err := server.Listen()
		Expect(err).NotTo(HaveOccurred())
		go func() {
			_, err := net.Dial("tcp", server.listener.Addr().String())
			Expect(err).NotTo(HaveOccurred())
			done <- true
		}()
		conn, ok := <-conns
		Expect(conn).NotTo(BeNil())
		Expect(ok).To(BeTrue())
		<-done
	})

	It("Errors on an invalid addr", func() {
		server.Addr = "localhost:65566"
		_, err := server.Listen()
		Expect(err).To(HaveOccurred())
	})
})
