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
	"strings"
	"time"
)

// ParseMessage creates a Message from a raw IRC message
func ParseMessage(raw string) *Message {
	message := &Message{}
	message.time = time.Now()
	// Message tags
	if strings.HasPrefix(raw, "@") {
		i := strings.IndexByte(raw, ' ')
		message.tags = parseTags(raw[1:i])
		raw = raw[i+1:]
	}
	// Message prefix
	if strings.HasPrefix(raw, ":") {
		i := strings.IndexByte(raw, ' ')
		message.prefix = raw[1:i]
		raw = raw[i+1:]
	}
	return message
}

func parseTags(tagString string) []Tag {
	tagStrings := strings.Split(tagString, ";")
	// http://ircv3.net/specs/core/message-tags-3.2.html#escaping-values
	r := strings.NewReplacer(
		`\:`, `;`,
		`\s`, ` `,
		`\\`, `\`,
		`\r`, "\r",
		`\n`, "\n",
	)
	tags := []Tag{}
	for i, tag := range tagStrings {
		kv := strings.SplitN(tag, "=", 2)
		tags[i].key = kv[0]
		tags[i].value = r.Replace(kv[1])
	}
	return tags
}

// Tag represents IRCv3 message tags, see:
// - http://ircv3.net/specs/core/message-tags-3.2.html
// - http://ircv3.net/specs/core/message-tags-3.3.html
//
// <tag> ::= <key> ['=' <escaped value>]
type Tag struct {
	// <key>           ::= [ <vendor> '/' ] <sequence of letters, digits, hyphens (`-`)>
	key string
	// <escaped value> ::= <sequence of any characters except NUL, CR, LF, semicolon
	//                     (`;`) and SPACE>
	// value is unescaped
	value string
}

// Message represents a parsed Message, see:
// - https://tools.ietf.org/html/rfc1459#section-2.3.1
//
// <message> ::= ['@' <tags> <SPACE>] [':' <prefix> <SPACE> ] <command> <params> <crlf>
type Message struct {
	// <tags>     ::= <tag> [';' <tag>]*
	tags []Tag
	// <prefix>   ::= <servername> | <nick> [ '!' <user> ] [ '@' <host> ]
	prefix string
	// <command>  ::= <letter> { <letter> } | <number> <number> <number>
	command string
	// <params>   ::= <SPACE> [ ':' <trailing> | <middle> <params> ]
	// <middle>   ::= <Any *non-empty* sequence of octets not including SPACE
	//                     or NUL or CR or LF, the first of which may not be ':'>
	// <trailing> ::= <Any, possibly *empty*, sequence of octets not including
	//                     NUL or CR or LF>
	params []string
	// time the message was received, may be replaced with server indicated time
	// if server-time is supported
	time time.Time
}
