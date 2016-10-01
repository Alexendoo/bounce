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

// IRCv3.1
const (
	// MultiPrefix - http://ircv3.net/specs/extensions/multi-prefix-3.1.html
	MultiPrefix = "multi-prefix"
	// Sasl --
	// - http://ircv3.net/specs/extensions/sasl-3.1.html
	// - http://ircv3.net/specs/extensions/sasl-3.2.html
	Sasl = "sasl"
	// AccountNotify - http://ircv3.net/specs/extensions/account-notify-3.1.html
	AccountNotify = "account-notify"
	// AwayNotify - http://ircv3.net/specs/extensions/away-notify-3.1.html
	AwayNotify = "away-notify"
	// ExtendedJoin - http://ircv3.net/specs/extensions/extended-join-3.1.html
	ExtendedJoin = "extended-join"
	// TLS - http://ircv3.net/specs/extensions/tls-3.1.html
	TLS = "tls"
)

// IRCv3.2
const (
	// Metadata - http://ircv3.net/specs/core/metadata-3.2.html
	Metadata = "metadata"
	// Monitor - http://ircv3.net/specs/core/monitor-3.2.html
	Monitor = "monitor"
	// AccountTag - http://ircv3.net/specs/extensions/account-tag-3.2.html
	AccountTag = "account-tag"
	// Batch - http://ircv3.net/specs/extensions/batch-3.2.html
	Batch = "batch"
	// CapNotify - http://ircv3.net/specs/extensions/cap-notify-3.2.html
	CapNotify = "cap-notify"
	// Chghost - http://ircv3.net/specs/extensions/chghost-3.2.html
	Chghost = "chghost"
	// EchoMessage - http://ircv3.net/specs/extensions/echo-message-3.2.html
	EchoMessage = "echo-message"
	// InviteNotify - http://ircv3.net/specs/extensions/invite-notify-3.2.html
	InviteNotify = "invite-notify"
	// ServerTime - http://ircv3.net/specs/extensions/server-time-3.2.html
	ServerTime = "server-time"
	// UserhostInNames - http://ircv3.net/specs/extensions/userhost-in-names-3.2.html
	UserhostInNames = "userhost-in-names"
)

// NewCapabilities returns a new Capabilities store
func NewCapabilities(supported map[string]string) *Capabilities {
	caps := &Capabilities{
		supported: make(map[string]string),
		enabled:   make(map[string]string),
	}
	if supported != nil {
		for cap, value := range supported {
			caps.supported[cap] = value
		}
	}
	return caps
}

// Capabilities represent a set of IRCv3 capabilites, with support for 3
// states: enabled, supported, not supported
//
// Safe for concurrent use
//
// - http://ircv3.net/specs/core/capability-negotiation-3.1.html
// - http://ircv3.net/specs/core/capability-negotiation-3.2.html
type Capabilities struct {
	sync.RWMutex

	// Version is the client capability negotiation version
	//
	// http://ircv3.net/specs/core/capability-negotiation-3.2.html#version-in-cap-ls
	Version int

	supported map[string]string
	enabled   map[string]string
}

// Supported returns if cap is supported
func (c *Capabilities) Supported(cap string) bool {
	c.RLock()
	_, supported := c.supported[cap]
	c.RUnlock()
	return supported
}

// SupportedValue returns the value of a supported cap
func (c *Capabilities) SupportedValue(cap string) string {
	c.RLock()
	value := c.supported[cap]
	c.RUnlock()
	return value
}

// Enabled returns if cap is enabled
func (c *Capabilities) Enabled(cap string) bool {
	c.RLock()
	_, enabled := c.enabled[cap]
	c.RUnlock()
	return enabled
}

// EnabledValue returns the value of an enabled cap
func (c *Capabilities) EnabledValue(cap string) string {
	c.RLock()
	value := c.enabled[cap]
	c.RUnlock()
	return value
}

// LS returns a map of the supported capabilites
//
// http://ircv3.net/specs/core/capability-negotiation-3.1.html#the-cap-ls-subcommand
func (c *Capabilities) LS() (supported map[string]string) {
	supported = make(map[string]string)
	c.RLock()
	for k, v := range c.supported {
		supported[k] = v
	}
	c.RUnlock()
	return supported
}

// List returns a map of the currently enabled capabilites
//
// http://ircv3.net/specs/core/capability-negotiation-3.1.html#the-cap-list-subcommand
func (c *Capabilities) List() (enabled map[string]string) {
	enabled = make(map[string]string)
	c.RLock()
	for k, v := range c.enabled {
		enabled[k] = v
	}
	c.RUnlock()
	return enabled
}

// Enable enables the provided caps, does NOT test if the capability is supported
func (c *Capabilities) Enable(caps map[string]string) {
	c.Lock()
	for cap, value := range caps {
		c.enabled[cap] = value
	}
	c.Unlock()
}

// Support marks the given caps as supported
//
// http://ircv3.net/specs/extensions/cap-notify-3.2.html#subcommands
func (c *Capabilities) Support(caps map[string]string) {
	c.Lock()
	for cap, value := range caps {
		c.supported[cap] = value
	}
	c.Unlock()
}

// Del removes support for and disables the given caps
//
// http://ircv3.net/specs/extensions/cap-notify-3.2.html#subcommands
func (c *Capabilities) Del(caps ...string) {
	c.Lock()
	for _, cap := range caps {
		delete(c.supported, cap)
		delete(c.enabled, cap)
	}
	c.Unlock()
}
