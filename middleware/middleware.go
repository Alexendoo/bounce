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

type UpstreamData struct {
	Message *irc.Message
	Client  *client.Client
	Network *network.Network
}

type DownstreamData struct {
	Message *irc.Message
	Clients []*client.Client
	Network *network.Network
}

// Middleware manipulates messages between the client[s] and network
type Middleware interface {
	upstream(data *UpstreamData, out chan<- *UpstreamData)
	downstream(data *DownstreamData, out chan<- *DownstreamData)
}

func NewUpstream(middleware ...Middleware) *Upstream {
	in := make(chan *UpstreamData)

	var out chan *UpstreamData
	for _, middleware := range middleware {
		if out == nil {
			out = pipeUpstream(middleware, in)
		} else {
			out = pipeUpstream(middleware, out)
		}
	}

	return &Upstream{In: in, Out: out}
}

func pipeUpstream(m Middleware, in chan *UpstreamData) chan *UpstreamData {
	out := make(chan *UpstreamData)
	go func() {
		for data := range in {
			m.upstream(data, out)
		}
		close(out)
	}()
	return out
}

type Upstream struct {
	In  chan<- *UpstreamData
	Out <-chan *UpstreamData
}

func NewDownstream(middleware ...Middleware) *Downstream {
	in := make(chan *DownstreamData)

	var out chan *DownstreamData
	for _, middleware := range middleware {
		if out == nil {
			out = pipeDownstream(middleware, in)
		} else {
			out = pipeDownstream(middleware, out)
		}
	}

	return &Downstream{In: in, Out: out}
}

func pipeDownstream(m Middleware, in chan *DownstreamData) chan *DownstreamData {
	out := make(chan *DownstreamData)
	go func() {
		for data := range in {
			m.downstream(data, out)
		}
		close(out)
	}()
	return out
}

type Downstream struct {
	In  chan<- *DownstreamData
	Out <-chan *DownstreamData
}
