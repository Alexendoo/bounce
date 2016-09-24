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
	"bufio"
	"fmt"
	"io"
	"net"
)

func Dial(network, address string) (Conn, error) {
	conn, err := net.Dial(network, address)
	return conn.(Conn), err
}

type Conn struct {
	net.Conn

	scanner *bufio.Scanner
}

func (conn *Conn) ReadLine() (line string, err error) {
	if conn.scanner == nil {
		conn.scanner = bufio.NewScanner(conn)
	}
	scanErr := conn.scanner.Scan()
	if scanErr {
		return "", conn.scanner.Err()
	}
	return conn.scanner.Text(), nil
}

func (conn *Conn) WriteLine(line string) (int, error) {
	return io.WriteString(conn, line+"\r\n")
}

func (conn *Conn) WriteLineF(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(conn, format+"\r\n", a)
}

// func Listen(network, laddr string) (Listener, error) {
// 	listener, err := net.Listen(network, laddr)
// 	return listener.(Listener), err
// }

// type Listener struct {
// 	net.Listener
// }
