package irc_test

import (
	"testing"

	. "macleod.io/bounce/irc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parses messages", func() {
	Specify("Command", func() {
		msg := ParseMessage("PRIVMSG")
		Expect(msg.Command).To(Equal("PRIVMSG"))
	})
	Specify("Tag", func() {
		msg := ParseMessage("@key=value PING")
		Expect(msg.Command).To(Equal("PING"))
		Expect(msg.Tags).To(HaveKeyWithValue("key", "value"))
		Expect(msg.Tags).To(HaveLen(1))
	})
	Specify("Multiple tags", func() {
		msg := ParseMessage("@first;second=2;third PING")
		Expect(msg.Command).To(Equal("PING"))
		Expect(msg.Tags).To(HaveLen(3))
	})
	Specify("Prefix", func() {
		msg := ParseMessage(":irc.example.org PING")
		Expect(msg.Command).To(Equal("PING"))
		Expect(msg.Prefix).To(Equal("irc.example.org"))
	})
	Specify("Param", func() {
		msg := ParseMessage("PING one")
		Expect(msg.Command).To(Equal("PING"))
		Expect(msg.Params).To(ConsistOf("one"))
	})
	Specify("Param with trailing whitespace", func() {
		msg := ParseMessage("PING one ")
		Expect(msg.Command).To(Equal("PING"))
		Expect(msg.Params).To(ConsistOf("one"))
	})
	Specify("Multiple params", func() {
		msg := ParseMessage("PRIVMSG #channel :a b c ")
		Expect(msg.Command).To(Equal("PRIVMSG"))
		Expect(msg.Params).To(BeEquivalentTo([]string{
			"#channel", "a b c ",
		}))
	})
	Specify("All of the above", func() {
		msg := ParseMessage(`@time=half\sfive;foo :example.org CAP * LS :server-time sasl`)
		Expect(msg.Tags).To(BeEquivalentTo(map[string]string{
			"time": "half five",
			"foo":  "",
		}))
		Expect(msg.Prefix).To(Equal("example.org"))
		Expect(msg.Command).To(Equal("CAP"))
		Expect(msg.Params).To(BeEquivalentTo([]string{
			"*", "LS", "server-time sasl",
		}))
	})
})

var _ = Describe("Composes messages", func() {
	Specify("Complex message", func() {
		msg := ParseMessage("@key=value :example.org PRIVMSG #channel :hello world ")
		Expect(msg.Buffer().String()).To(Equal("@key=value :example.org PRIVMSG #channel :hello world \r\n"))
	})
	Specify("Multiple tags", func() {
		msg := ParseMessage("@one=1;two PING")
		Expect(msg.Buffer().String()).To(SatisfyAny(
			Equal("@one=1;two PING\r\n"),
			Equal("@two;one=1 PING\r\n"),
		))
	})
})

func BenchmarkParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseMessage("@time=now :example.org PRIVMSG #channel :message message message")
	}
	b.StopTimer()
}

func BenchmarkBufferComposition(b *testing.B) {
	msg := &Message{
		Tags: map[string]string{
			"key": "value",
		},
		Prefix:  "example.org",
		Command: "PRIVMSG",
		Params: []string{
			"#channel",
			"hello world",
		},
	}
	for i := 0; i < b.N; i++ {
		msg.Buffer()
	}
	b.StopTimer()
}
