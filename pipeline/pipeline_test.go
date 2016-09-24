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

import (
	"math/rand"
	"testing"
	"time"

	c "github.com/smartystreets/goconvey/convey"
)

func TestPipeline(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	c.Convey("it should pass data", t, func() {
		pipeline := Plumb(waiter{})
		message := data{name: "one", value: 2}
		go func() {
			pipeline.Push(message)
			pipeline.Close()
		}()
		pipeline.Pull(func(out data) {
			c.So(out, c.ShouldResemble, message)
		})
	})

	c.Convey("it should retain order", t, func() {
		pipeline := Plumb(waiter{}, waiter{}, waiter{})
		messages := []data{
			{name: "one", value: 1},
			{name: "two", value: 2},
			{name: "three", value: 3},
		}
		go func() {
			for _, object := range messages {
				pipeline.Push(object)
			}
			pipeline.Close()
		}()
		i := 0
		pipeline.Pull(func(out data) {
			c.So(out, c.ShouldResemble, messages[i])
			i++
		})
	})
}

type waiter struct{}

func (w waiter) Process(in chan data) chan data {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return in
}
