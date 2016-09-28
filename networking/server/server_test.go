package network

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var server *Server

	BeforeEach(func() {
		server = &Server{
			Addr: "localhost:0",
		}
	})

	It("Accepts connections", func() {

	})
})
