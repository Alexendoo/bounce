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
	"io"
	"net"

	"macleod.io/bounce/irc"
)

type Server struct {
	Host string
	Port string

	listener net.Listener
}

func (s *Server) Listen() error {
	addr := net.JoinHostPort(s.Host, s.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.registerClient(conn)
	}
}

func (s *Server) registerClient(conn net.Conn) {
	io.WriteString(conn, "hello\n")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := irc.ParseMessage(scanner.Text())
		switch msg.Command {
		case "CAP":

		}
	}
	conn.Close()
}

func (s *Server) Close() error {
	return s.listener.Close()
}
