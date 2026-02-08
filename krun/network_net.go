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
func (c *Context) AddNetUnixStream(cfg NetUnixConfig) error {
	var cPath *C.char
	if cfg.Path != "" {
		cPath = C.CString(cfg.Path)
		defer C.free(unsafe.Pointer(cPath))
	}
	return checkRet(
		C.krun_add_net_unixstream(
			C.uint32_t(c.id), cPath, C.int(cfg.FD),
			(*C.uint8_t)(unsafe.Pointer(&cfg.MAC[0])),
			C.uint32_t(cfg.Features), C.uint32_t(cfg.Flags),
		),
		"krun_add_net_unixstream",
	)
}

// AddNetUnixGram adds a virtio-net device with a unixgram-based backend
// (e.g., gvproxy or vmnet-helper).
// If using gvproxy in vfkit mode with a path, include [NetFlagVfkit] in Flags.
func (c *Context) AddNetUnixGram(cfg NetUnixConfig) error {
	var cPath *C.char
	if cfg.Path != "" {
		cPath = C.CString(cfg.Path)
		defer C.free(unsafe.Pointer(cPath))
	}
	return checkRet(
		C.krun_add_net_unixgram(
			C.uint32_t(c.id), cPath, C.int(cfg.FD),
			(*C.uint8_t)(unsafe.Pointer(&cfg.MAC[0])),
			C.uint32_t(cfg.Features), C.uint32_t(cfg.Flags),
		),
		"krun_add_net_unixgram",
	)
}

// AddNetTap adds a virtio-net device with the TAP backend.
func (c *Context) AddNetTap(cfg NetTapConfig) error {
	cTapName := C.CString(cfg.TapName)
	defer C.free(unsafe.Pointer(cTapName))
	return checkRet(
		C.krun_add_net_tap(
			C.uint32_t(c.id), cTapName,
			(*C.uint8_t)(unsafe.Pointer(&cfg.MAC[0])),
			C.uint32_t(cfg.Features), C.uint32_t(cfg.Flags),
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
