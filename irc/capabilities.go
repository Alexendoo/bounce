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

import "sync"

func NewCapabilities() *Capabilities {
	return &Capabilities{
		caps: make(map[string]bool),
	}
}

// Capabilities
type Capabilities struct {
	sync.RWMutex
	caps map[string]bool
}

func (c *Capabilities) Supported(cap string) bool {
	c.RLock()
	_, exists := c.caps[cap]
	c.RUnlock()

	return exists
}

func (c *Capabilities) Enabled(cap string) bool {
	c.RLock()
	enabled := c.caps[cap]
	c.RUnlock()

	return enabled
}

func (c *Capabilities) Support(caps ...string) {
	c.Lock()
	for _, cap := range caps {
		_, exists := c.caps[cap]
		if !exists {
			c.caps[cap] = false
		}
	}
	c.Unlock()
}

func (c *Capabilities) Enable(caps ...string) {
	c.Lock()
	for _, cap := range caps {
		c.caps[cap] = true
	}
	c.Unlock()
}

func (c *Capabilities) Disable(caps ...string) {
	c.Lock()
	for _, cap := range caps {
		c.caps[cap] = false
	}
	c.Unlock()
}
