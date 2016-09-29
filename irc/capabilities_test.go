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

package irc_test

import (
	. "macleod.io/bounce/irc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Capabilities", func() {
	var caps *Capabilities

	BeforeEach(func() {
		caps = NewCapabilities()
	})

	It("Enables Capabilities", func() {
		caps.Enable("server-time", "batch")
		Expect(caps.Enabled("server-time")).To(BeTrue())
		Expect(caps.Enabled("batch")).To(BeTrue())
		Expect(caps.Enabled("away-notify")).To(BeFalse())
	})

	It("Disables Capabilities", func() {
		caps.Enable("server-time", "batch")
		caps.Disable("server-time", "batch")
		Expect(caps.Enabled("server-time")).To(BeFalse())
		Expect(caps.Enabled("batch")).To(BeFalse())
	})

	It("Distinguishes supported and enabled", func() {
		caps.Support("server-time")

		Expect(caps.Supported("batch")).To(BeFalse())
		Expect(caps.Supported("server-time")).To(BeTrue())
		Expect(caps.Enabled("server-time")).To(BeFalse())
	})

	It("Lists enabled caps", func() {
		caps.Support("batch", "away-notify")
		caps.Enable("server-time", "sasl")

		Expect(caps.List()).To(ConsistOf("server-time", "sasl"))
	})
})
