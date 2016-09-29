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
	n.In = make(chan *irc.Message)
	n.Out = make(chan *irc.Message)
	conn, err := net.Dial("tcp", n.Addr)
	if err != nil {
		return err
	}
	n.conn = conn
	go n.accept()
	go n.scan()
	return nil
}

func (n *Network) accept() {
	for message := range n.In {
		message.Buffer().WriteTo(n.conn)
	}
}

func (n *Network) scan() {
	scanner := bufio.NewScanner(n.conn)
	for scanner.Scan() {
		n.Out <- irc.ParseMessage(scanner.Text())
	}
	// TODO: Reconnect + emit error to clients
}

func (n *Network) Register() error {
	_, err := fmt.Fprintf(
		n.conn,
		"CAP LS 302\r\n"+
			"NICK %s\r\n"+
			"USER %s - - :%s\r\n",
		n.Nick, n.User, n.Real,
	)
	if err != nil {
		return err
	}
	return nil
}

func (n *Network) Close() error {
	close(n.In)
	return n.conn.Close()
}
