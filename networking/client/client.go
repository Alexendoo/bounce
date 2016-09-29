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
	"net"

	"macleod.io/bounce/irc"
)

func New(conn net.Conn) *Client {
	client := &Client{
		conn: conn,
		In:   make(chan *irc.Message),
		Out:  make(chan *irc.Message),
	}
	go client.accept()
	go client.scan()
	return client
}

type Capabilities struct {
	Available    []string
	Acknowledged []string
}

type Client struct {
	Capabilities []Capabilities
	conn         net.Conn

	In  chan *irc.Message
	Out chan *irc.Message
}

func (c *Client) Close() error {
	close(c.In)
	return c.conn.Close()
}

func (c *Client) accept() {
	for message := range c.In {
		// TODO : middleware
		message.Buffer().WriteTo(c.conn)
	}
}

func (c *Client) scan() {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		message := irc.ParseMessage(scanner.Text())
		// TODO : middleware
		c.Out <- message
	}
	close(c.Out)
}
