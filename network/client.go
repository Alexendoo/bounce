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
	"container/list"
	"io"
	"net"
)

// TODO : register per client, hook into main channel? draw graph
//      : capabilities under package irc

type Capabilities struct {
	Available    []string
	Acknowledged []string
}

type Client struct {
	Server       *Server
	Element      *list.Element
	Capabilities []Capabilities
	Conn         net.Conn
}

func (c *Client) Register(clients *list.List) {
	io.WriteString(c.Conn, "hello\n")
}
