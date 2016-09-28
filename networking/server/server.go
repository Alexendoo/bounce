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

import "net"

type Server struct {
	Addr string

	listener net.Listener
}

func (s *Server) Listen() (chan net.Conn, error) {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return nil, err
	}
	out := make(chan net.Conn)
	s.listener = listener
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				return
			}
			out <- conn
		}
	}()
	return out, nil
}

func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}

	return s.listener.Close()
}
