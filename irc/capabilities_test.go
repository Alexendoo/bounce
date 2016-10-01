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
	var initialValues map[string]string

	BeforeEach(func() {
		initialValues = map[string]string{
			ServerTime: "",
			Sasl:       "PLAIN",
		}
		caps = NewCapabilities(initialValues)
	})

	It("Detects supported caps", func() {
		Expect(caps.Supported(ServerTime)).To(BeTrue())
	})

	It("Detects unsupported caps", func() {
		Expect(caps.Supported(AwayNotify)).To(BeFalse())
	})

	It("Enables new caps", func() {
		caps.Enable(map[string]string{
			Sasl: "PLAIN",
		})
		Expect(caps.Enabled(Sasl)).To(BeTrue())
		Expect(caps.Enabled(ServerTime)).To(BeFalse())
		Expect(caps.EnabledValue(Sasl)).To(Equal("PLAIN"))
	})

	It("Supports new caps", func() {
		caps.Support(map[string]string{
			Batch: "",
			"key": "value",
		})

		Expect(caps.Supported(Batch)).To(BeTrue())
		Expect(caps.Supported("key")).To(BeTrue())
		Expect(caps.SupportedValue("key")).To(Equal("value"))
	})

	It("Deletes caps", func() {
		caps.Enable(map[string]string{
			ServerTime: "",
		})

		caps.Del(ServerTime, Sasl)
		Expect(caps.LS()).To(BeEmpty())
		Expect(caps.List()).To(BeEmpty())
	})

	It("Lists supported capabilities", func() {
		Expect(caps.LS()).To(Equal(initialValues))
	})

	It("Lists enabled capabilities", func() {
		initial := map[string]string{ServerTime: ""}

		caps.Enable(initial)
		list := caps.List()
		Expect(list).To(Equal(initial))
	})
})
