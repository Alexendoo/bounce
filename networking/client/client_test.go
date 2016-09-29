package client_test

import (
	"bufio"
	"net"

	"macleod.io/bounce/irc"
	. "macleod.io/bounce/networking/client"

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
