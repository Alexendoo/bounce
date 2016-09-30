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

package middleware

import (
	"macleod.io/bounce/irc"
	"macleod.io/bounce/networking/client"
	"macleod.io/bounce/networking/network"
)

type UpstreamMessage struct {
	Message *irc.Message
	Client  *client.Client
	Network *network.Network
}

type DownstreamMessage struct {
	Message *irc.Message
	Clients []*client.Client
	Network *network.Network
}

// Middleware manipulates messages between the client[s] and network
type Middleware interface {
	upstream(message *UpstreamMessage, out chan<- *UpstreamMessage)
	downstream(message *DownstreamMessage, out chan<- *DownstreamMessage)
}

func NewUpstream(middleware ...Middleware) *Upstream {
	in := make(chan *UpstreamMessage)

	var out chan *UpstreamMessage
	for _, middleware := range middleware {
		if out == nil {
			out = pipeUpstream(middleware, in)
		} else {
			out = pipeUpstream(middleware, out)
		}
	}

	return &Upstream{In: in, Out: out}
}

func pipeUpstream(m Middleware, in chan *UpstreamMessage) chan *UpstreamMessage {
	out := make(chan *UpstreamMessage)
	go func() {
		for message := range in {
			m.upstream(message, out)
		}
		close(out)
	}()
	return out
}

type Upstream struct {
	In  chan<- *UpstreamMessage
	Out <-chan *UpstreamMessage
}

func NewDownstream(middleware ...Middleware) *Downstream {
	in := make(chan *DownstreamMessage)

	var out chan *DownstreamMessage
	for _, middleware := range middleware {
		if out == nil {
			out = pipeDownstream(middleware, in)
		} else {
			out = pipeDownstream(middleware, out)
		}
	}

	return &Downstream{In: in, Out: out}
}

func pipeDownstream(m Middleware, in chan *DownstreamMessage) chan *DownstreamMessage {
	out := make(chan *DownstreamMessage)
	go func() {
		for message := range in {
			m.downstream(message, out)
		}
		close(out)
	}()
	return out
}

type Downstream struct {
	In  chan<- *DownstreamMessage
	Out <-chan *DownstreamMessage
}
