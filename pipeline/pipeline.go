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

package pipeline

type data struct {
	value int
	name  string
}

type Pipe interface {
	Process(in chan data) chan data
}

type Pipeline struct {
	upstream   chan data
	downstream chan data
}

func (p *Pipeline) Push(item data) {
	p.upstream <- item
}

func (p *Pipeline) Pull(handler func(data)) {
	for i := range p.downstream {
		handler(i)
	}
}

func (p *Pipeline) Close() {
	close(p.upstream)
}

func Plumb(pipes ...Pipe) Pipeline {
	upstream := make(chan data)
	var downstream chan data
	for _, pipe := range pipes {
		if downstream == nil {
			downstream = pipe.Process(upstream)
		} else {
			downstream = pipe.Process(downstream)
		}
	}
	return Pipeline{upstream: upstream, downstream: downstream}
}
