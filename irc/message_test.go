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

package irc

import (
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func TestParsing(t *testing.T) {
	c.Convey("basic command", t, func() {
		msg := ParseMessage("PRIVMSG")
		c.So(msg, c.ShouldResemble, &Message{
			Tags:    map[string]string(nil),
			Prefix:  "",
			Command: "PRIVMSG",
			Params:  []string(nil),
			Time:    msg.Time,
		})
	})
	c.Convey("basic tag", t, func() {
		msg := ParseMessage("@key=value PING")
		c.So(msg, c.ShouldResemble, &Message{
			Tags: map[string]string{
				"key": "value",
			},
			Prefix:  "",
			Command: "PING",
			Params:  []string(nil),
			Time:    msg.Time,
		})
	})
	c.Convey("multiple tags", t, func() {
		msg := ParseMessage("@first;second=2;third PING")
		c.So(msg, c.ShouldResemble, &Message{
			Tags: map[string]string{
				"first":  "",
				"second": "2",
				"third":  "",
			},
			Prefix:  "",
			Command: "PING",
			Params:  []string(nil),
			Time:    msg.Time,
		})
	})
	c.Convey("prefix", t, func() {
		msg := ParseMessage(":irc.example.org PING")
		c.So(msg, c.ShouldResemble, &Message{
			Tags:    map[string]string(nil),
			Prefix:  "irc.example.org",
			Command: "PING",
			Params:  []string(nil),
			Time:    msg.Time,
		})
	})
	c.Convey("basic param", t, func() {
		msg := ParseMessage("PING one")
		c.So(msg, c.ShouldResemble, &Message{
			Tags:    map[string]string(nil),
			Prefix:  "",
			Command: "PING",
			Params:  []string{"one"},
			Time:    msg.Time,
		})
	})
	c.Convey("trailing whitespace", t, func() {
		msg := ParseMessage("PING ")
		c.So(msg, c.ShouldResemble, &Message{
			Tags:    map[string]string(nil),
			Prefix:  "",
			Command: "PING",
			Params:  []string(nil),
			Time:    msg.Time,
		})
	})
	c.Convey("trailing whitespace param", t, func() {
		msg := ParseMessage("PING one ")
		c.So(msg, c.ShouldResemble, &Message{
			Tags:    map[string]string(nil),
			Prefix:  "",
			Command: "PING",
			Params:  []string{"one"},
			Time:    msg.Time,
		})
	})
	c.Convey("multiple params", t, func() {
		msg := ParseMessage("PRIVMSG #channel :a b c ")
		c.So(msg, c.ShouldResemble, &Message{
			Tags:    map[string]string(nil),
			Prefix:  "",
			Command: "PRIVMSG",
			Params: []string{
				"#channel",
				"a b c ",
			},
			Time: msg.Time,
		})
	})
	c.Convey("all of the above", t, func() {
		msg := ParseMessage(`@time=half\sfive;foo :example.org CAP * LS :server-time sasl`)
		c.So(msg, c.ShouldResemble, &Message{
			Tags: map[string]string{
				"time": "half five",
				"foo":  "",
			},
			Prefix:  "example.org",
			Command: "CAP",
			Params: []string{
				"*",
				"LS",
				"server-time sasl",
			},
			Time: msg.Time,
		})
	})
}

func TestNextToken(t *testing.T) {
	c.Convey("initial offset", t, func() {
		trail, lead := nextToken("ABC DEF GHI", 0, 0)
		c.So(trail, c.ShouldEqual, 0)
		c.So(lead, c.ShouldEqual, 3)
	})
	c.Convey("middle offset", t, func() {
		trail, lead := nextToken("ABC DEF GHI", 0, 3)
		c.So(trail, c.ShouldEqual, 4)
		c.So(lead, c.ShouldEqual, 7)
	})
	c.Convey("penultimate offset", t, func() {
		trail, lead := nextToken("ABC DEF GHI", 4, 7)
		c.So(trail, c.ShouldEqual, 8)
		c.So(lead, c.ShouldEqual, 11)
	})
	c.Convey("last offset", t, func() {
		trail, lead := nextToken("ABC DEF GHI", 8, 11)
		c.So(trail, c.ShouldEqual, 11)
		c.So(lead, c.ShouldEqual, 11)
	})
	c.Convey("skips multiple spaces", t, func() {
		trail, lead := nextToken("ABC  DEF", 0, 4)
		c.So(trail, c.ShouldEqual, 5)
		c.So(lead, c.ShouldEqual, 8)
	})
}

func BenchmarkParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseMessage("@time=now :example.org PRIVMSG #channel :message message message")
	}
	b.StopTimer()
}

func BenchmarkNextToken(b *testing.B) {
	for i := 0; i < b.N; i += 4 {
		nextToken("ABC DEF    GHI", 0, 0)
		nextToken("ABC DEF    GHI", 0, 3)
		nextToken("ABC DEF    GHI", 4, 7)
		nextToken("ABC DEF    GHI", 11, 14)
	}
	b.StopTimer()
}
