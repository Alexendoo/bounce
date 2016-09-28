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

package network

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"macleod.io/bounce/irc"
)

type Network struct {
	Addr string

	Nick string
	Real string
	User string

	In  chan *irc.Message
	Out chan *irc.Message

	conn net.Conn
}

func (n *Network) Connect() error {
	conn, err := net.Dial("tcp", n.Addr)
	if err != nil {
		return err
	}
	n.conn = conn
	go n.scan()
	go n.accept()
	return nil
}

func (n *Network) scan() {
	scanner := bufio.NewScanner(n.conn)
	for scanner.Scan() {
		n.Out <- irc.ParseMessage(scanner.Text())
	}
}

func (n *Network) accept() {
	for message := range n.In {
		message.Buffer().WriteTo(n.conn)
	}
}

func (n *Network) register() {
	n.sendRaw("CAP LS 302")
	n.sendRaw(fmt.Sprintf("NICK %s", n.Nick))
	n.sendRaw(fmt.Sprintf("USER %s - - :%s", n.User, n.Real))
}

func (n *Network) sendRaw(message string) {
	io.WriteString(n.conn, message+"\r\n")
}

func (n *Network) Disconnect() {
	n.conn.Close()
}
