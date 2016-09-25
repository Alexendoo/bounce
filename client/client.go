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
	"fmt"
	"io"
	"net"

	"macleod.io/bounce/irc"
)

type Network struct {
	Name string

	Host string
	Port string

	Nick string
	Real string
	User string

	conn net.Conn
}

func (n *Network) connect() error {
	addr := net.JoinHostPort(n.Host, n.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	n.conn = conn
	return nil
}

func (n *Network) disconnect() {
	n.conn.Close()
}

func (n *Network) register() {
	n.sendRaw("CAP LS 302")
	n.sendRaw(fmt.Sprintf("NICK %s", n.Nick))
	n.sendRaw(fmt.Sprintf("USER %s - - :%s", n.User, n.Real))
}

func (n *Network) sendRaw(message string) {
	io.WriteString(n.conn, message+"\r\n")
}

func (n *Network) listen() chan *irc.Message {
	messages := make(chan *irc.Message)
	scanner := bufio.NewScanner(n.conn)
	go func() {
		for scanner.Scan() {
			messages <- irc.ParseMessage(scanner.Text())
		}
		close(messages)
	}()
	return messages
}
