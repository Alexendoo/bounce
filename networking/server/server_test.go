package server

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
