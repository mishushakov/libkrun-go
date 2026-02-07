package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// AddVsockPort maps a vsock port to a host UNIX socket path.
func (c *Context) AddVsockPort(port uint32, filepath string) error {
	cPath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(
		C.krun_add_vsock_port(C.uint32_t(c.id), C.uint32_t(port), cPath),
		"krun_add_vsock_port",
	)
}

// AddVsockPort2 maps a vsock port to a host UNIX socket path with a listen mode option.
// If listen is true, the guest expects connections to be initiated from the host side.
func (c *Context) AddVsockPort2(port uint32, filepath string, listen bool) error {
	cPath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(
		C.krun_add_vsock_port2(C.uint32_t(c.id), C.uint32_t(port), cPath, C.bool(listen)),
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
