package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// VsockPortConfig configures a vsock port mapping.
type VsockPortConfig struct {
	Port   uint32
	Path   string
	Listen bool // false = guest initiates connections
}

// AddVsockPort maps a vsock port to a host UNIX socket path.
func (c *Context) AddVsockPort(cfg VsockPortConfig) error {
	cPath := C.CString(cfg.Path)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(
		C.krun_add_vsock_port2(C.uint32_t(c.id), C.uint32_t(cfg.Port), cPath, C.bool(cfg.Listen)),
		"krun_add_vsock_port2",
	)
}

// AddVsock adds a vsock device with specified TSI features.
// Call [Context.DisableImplicitVsock] before using this to disable the default vsock.
// tsiFeatures is a bitmask of TSIHijack* flags. Use 0 for no TSI hijacking.
func (c *Context) AddVsock(tsiFeatures uint32) error {
	return checkRet(
		C.krun_add_vsock(C.uint32_t(c.id), C.uint32_t(tsiFeatures)),
		"krun_add_vsock",
	)
}

// DisableImplicitVsock disables the automatically created vsock device.
func (c *Context) DisableImplicitVsock() error {
	return checkRet(
		C.krun_disable_implicit_vsock(C.uint32_t(c.id)),
		"krun_disable_implicit_vsock",
	)
}
