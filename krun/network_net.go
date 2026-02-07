//go:build krun_net

package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// AddNetUnixStream adds a virtio-net device connected to a unixstream-based
// network proxy (e.g., passt or socket_vmnet).
//
// path and fd are mutually exclusive: pass "" for path when using fd,
// or -1 for fd when using path.
// mac is the 6-byte MAC address for the interface.
// features is a bitmask of NetFeature* flags.
// flags is a bitmask of NetFlag* flags.
func (c *Context) AddNetUnixStream(path string, fd int, mac [6]byte, features, flags uint32) error {
	var cPath *C.char
	if path != "" {
		cPath = C.CString(path)
		defer C.free(unsafe.Pointer(cPath))
	}
	return checkRet(
		C.krun_add_net_unixstream(
			C.uint32_t(c.id), cPath, C.int(fd),
			(*C.uint8_t)(unsafe.Pointer(&mac[0])),
			C.uint32_t(features), C.uint32_t(flags),
		),
		"krun_add_net_unixstream",
	)
}

// AddNetUnixGram adds a virtio-net device with a unixgram-based backend
// (e.g., gvproxy or vmnet-helper).
//
// path and fd are mutually exclusive: pass "" for path when using fd,
// or -1 for fd when using path.
// mac is the 6-byte MAC address for the interface.
// features is a bitmask of NetFeature* flags.
// flags is a bitmask of NetFlag* flags.
// If using gvproxy in vfkit mode with path, include [NetFlagVfkit] in flags.
func (c *Context) AddNetUnixGram(path string, fd int, mac [6]byte, features, flags uint32) error {
	var cPath *C.char
	if path != "" {
		cPath = C.CString(path)
		defer C.free(unsafe.Pointer(cPath))
	}
	return checkRet(
		C.krun_add_net_unixgram(
			C.uint32_t(c.id), cPath, C.int(fd),
			(*C.uint8_t)(unsafe.Pointer(&mac[0])),
			C.uint32_t(features), C.uint32_t(flags),
		),
		"krun_add_net_unixgram",
	)
}

// AddNetTap adds a virtio-net device with the TAP backend.
// tapName is the TAP device name.
// mac is the 6-byte MAC address for the interface.
func (c *Context) AddNetTap(tapName string, mac [6]byte, features, flags uint32) error {
	cTapName := C.CString(tapName)
	defer C.free(unsafe.Pointer(cTapName))
	return checkRet(
		C.krun_add_net_tap(
			C.uint32_t(c.id), cTapName,
			(*C.uint8_t)(unsafe.Pointer(&mac[0])),
			C.uint32_t(features), C.uint32_t(flags),
		),
		"krun_add_net_tap",
	)
}

// SetNetMac sets the MAC address for the virtio-net device when using the
// passt backend.
func (c *Context) SetNetMac(mac [6]byte) error {
	return checkRet(
		C.krun_set_net_mac(C.uint32_t(c.id), (*C.uint8_t)(unsafe.Pointer(&mac[0]))),
		"krun_set_net_mac",
	)
}
