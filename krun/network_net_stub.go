//go:build !krun_net

package krun

import "syscall"

// AddNetUnixStream adds a virtio-net device connected to a unixstream-based network proxy.
// Requires building with -tags krun_net.
func (c *Context) AddNetUnixStream(cfg NetUnixConfig) error {
	return &Error{Func: "krun_add_net_unixstream", Errno: syscall.ENOSYS}
}

// AddNetUnixGram adds a virtio-net device with a unixgram-based backend.
// Requires building with -tags krun_net.
func (c *Context) AddNetUnixGram(cfg NetUnixConfig) error {
	return &Error{Func: "krun_add_net_unixgram", Errno: syscall.ENOSYS}
}

// AddNetTap adds a virtio-net device with the TAP backend.
// Requires building with -tags krun_net.
func (c *Context) AddNetTap(cfg NetTapConfig) error {
	return &Error{Func: "krun_add_net_tap", Errno: syscall.ENOSYS}
}

// SetNetMac sets the MAC address for the virtio-net device.
// Requires building with -tags krun_net.
func (c *Context) SetNetMac(mac [6]byte) error {
	return &Error{Func: "krun_set_net_mac", Errno: syscall.ENOSYS}
}
