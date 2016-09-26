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
	"bytes"
	"strings"
	"time"
)

// ParseMessage creates a Message from a raw IRC message without the <crlf>
//
// <message> ::= ['@' <tags> <SPACE>] [':' <prefix> <SPACE> ] <command> <params>
func ParseMessage(raw string) *Message {
	message := &Message{}
	message.Time = time.Now()
	// <tags> ::= <tag> [';' <tag>]*
	trail, lead := nextToken(raw, 0, 0)
	if raw[trail] == '@' {
		message.Tags = parseTags(raw[trail+1 : lead])
		trail, lead = nextToken(raw, trail, lead)
	}
	// <prefix> ::= <servername> | <nick> [ '!' <user> ] [ '@' <host> ]
	if raw[trail] == ':' {
		message.Prefix = raw[trail+1 : lead]
		trail, lead = nextToken(raw, trail, lead)
	}
	// <command>  ::= <letter> { <letter> } | <number> <number> <number>
	message.Command = raw[trail:lead]
	trail, lead = nextToken(raw, trail, lead)
	// <params> ::= <SPACE> [ ':' <trailing> | <middle> <params> ]
	length := len(raw)
	for trail < length {
		// <trailing> ::= <Any, possibly *empty*, sequence of octets not including
		//                 NUL or CR or LF>
		if raw[trail] == ':' {
			message.Params = append(message.Params, raw[trail+1:])
			break
		}
		// <middle> ::= <Any *non-empty* sequence of octets not including SPACE or
		//               NUL or CR or LF, the first of which may not be ':'>
		message.Params = append(message.Params, raw[trail:lead])
		trail, lead = nextToken(raw, trail, lead)
	}
	return message
}

// advance the trail and lead cursors to the start and end of the next word
//
// e.g. nextToken("ABC DEF", 0, 3) = (4, 7)
// "ABC DEF" → "ABC DEF"
//  ^  ^     →      ^  ^
func nextToken(raw string, trail, lead int) (newTrail, newLead int) {
	length := len(raw)
	// advance trail cursor to the beginning of the next word
	for trail = lead; trail < length; trail++ {
		if raw[trail] != ' ' {
			break
		}
	}
	// advance lead cursor to the end of the next word
	for lead = trail; lead < length; lead++ {
		if raw[lead] == ' ' {
			break
		}
	}
	return trail, lead
}

func parseTags(tagString string) map[string]string {
	// <tags> ::= <tag> [';' <tag>]*
	tagStrings := strings.Split(tagString, ";")
	// http://ircv3.net/specs/core/message-tags-3.2.html#escaping-values
	r := strings.NewReplacer(
		`\:`, `;`,
		`\s`, ` `,
		`\\`, `\`,
		`\r`, "\r",
		`\n`, "\n",
	)
	tags := make(map[string]string)
	for _, tag := range tagStrings {
		// <tag>           ::= <key> ['=' <escaped value>]
		// <key>           ::= [ <vendor> '/' ] <sequence of letters, digits,
		//                                       hyphens (`-`)>
		// <escaped value> ::= <sequence of any characters except NUL, CR, LF,
		//                      semicolon (`;`) and SPACE>
		kv := strings.SplitN(tag, "=", 2)
		if len(kv) == 2 {
			tags[kv[0]] = r.Replace(kv[1])
		} else {
			tags[tag] = ""
		}
	}
	return tags
}

// Message represents a parsed Message, see:
// - https://tools.ietf.org/html/rfc1459#section-2.3
// - http://ircv3.net/specs/core/message-tags-3.2.html
type Message struct {
	Tags map[string]string
	// Prefix is typically the source of the message, e.g. nick!ident@host
	// or irc.example.org
	Prefix string
	// Command is a regular command (PRIVMSG, PING, etc.) or a numeric reply
	// - https://tools.ietf.org/html/rfc1459#section-4
	// - https://tools.ietf.org/html/rfc1459#section-6
	Command string
	// Params depend on the Command, for example with a PRIVMSG it would be the
	// message target followed by the message text to send
	Params []string
	// Time is the message was received
	Time time.Time
}

// Buffer returns a buffer containing the string form of the Message,
// including crlf
//
// ['@' <tags> <SPACE>] [':' <prefix> <SPACE> ] <command> <params> <crlf>
func (m *Message) Buffer() *bytes.Buffer {
	var buffer bytes.Buffer
	// <tags> ::= <tag> [';' <tag>]*
	if len(m.Tags) > 0 {
		buffer.WriteByte('@')
		appendTags(&buffer, m)
		buffer.WriteByte(' ')
	}
	// <prefix> ::= <servername> | <nick> [ '!' <user> ] [ '@' <host> ]
	if len(m.Prefix) > 0 {
		buffer.WriteByte(':')
		buffer.WriteString(m.Prefix)
		buffer.WriteByte(' ')
	}
	// <command>  ::= <letter> { <letter> } | <number> <number> <number>
	buffer.WriteString(m.Command)
	// <params> ::= <SPACE> [ ':' <trailing> | <middle> <params> ]
	for _, param := range m.Params {
		buffer.WriteByte(' ')
		if strings.ContainsRune(param, ' ') {
			buffer.WriteByte(':')
		}
		// <trailing> ::= <Any, possibly *empty*, sequence of octets not including
		//                 NUL or CR or LF>
		// <middle>   ::= <Any *non-empty* sequence of octets not including SPACE
		//                 or NUL or CR or LF, the first of which may not be ':'>
		buffer.WriteString(param)
	}
	buffer.WriteString("\r\n")
	return &buffer
}

// String returns the string form of the Message including the crlf
//
// <message> ::= ['@' <tags> <SPACE>] [':' <prefix> <SPACE> ] <command> <params> <crlf>
func (m *Message) String() string {
	return m.Buffer().String()
}

// appendTags writes the string representation of m.Tags to buffer
func appendTags(buffer *bytes.Buffer, m *Message) {
	var multiple bool
	r := strings.NewReplacer(
		`;`, `\:`,
		` `, `\s`,
		`\`, `\\`,
		"\r", `\r`,
		"\n", `\n`,
	)
	// <tags> ::= <tag> [';' <tag>]*
	for key, value := range m.Tags {
		// <tag> ::= <key> ['=' <escaped value>]
		if multiple {
			buffer.WriteByte(';')
		}
		// <key> ::= [ <vendor> '/' ] <sequence of letters, digits, hyphens (`-`)>
		buffer.WriteString(key)
		// <escaped value> ::= <sequence of any characters except NUL, CR, LF,
		//                      semicolon (`;`) and SPACE>
		if value != "" {
			buffer.WriteByte('=')
			buffer.WriteString(r.Replace(value))
		}
		multiple = true
	}

}
