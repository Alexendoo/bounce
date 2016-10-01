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

package middleware_test

import (
	"macleod.io/bounce/irc"
	. "macleod.io/bounce/middleware"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Middleware", func() {
	var message *irc.Message

	BeforeEach(func() {
		message = irc.ParseMessage(":example.org PING :v2FaU")
	})

	Context("Upstream", func() {
		var upstream *Upstream

		BeforeEach(func() {
			upstream = NewUpstream(&Null{}, &Null{}, &Null{})
		})

		It("Passes messages", func() {
			data := &UpstreamData{
				Message: message,
			}
			upstream.In <- data
			Expect(<-upstream.Out).To(Equal(data))
		})

		It("Closes", func() {
			close(upstream.In)
			_, ok := <-upstream.Out
			Expect(ok).To(BeFalse())
		})
	})

	Context("Downstream", func() {
		var downstream *Downstream

		BeforeEach(func() {
			downstream = NewDownstream(&Null{}, &Null{}, &Null{})
		})

		It("Passes messages", func() {
			data := &DownstreamData{
				Message: message,
			}
			downstream.In <- data
			Expect(<-downstream.Out).To(Equal(data))
		})

		It("Closes", func() {
			close(downstream.In)
			_, ok := <-downstream.Out
			Expect(ok).To(BeFalse())
		})
	})
})
