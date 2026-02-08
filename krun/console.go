package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// SetConsoleOutput redirects the implicit console output to a file.
// This only applies to the implicitly created console and has no effect
// if the implicit console is disabled via [Context.DisableImplicitConsole].
func (c *Context) SetConsoleOutput(filepath string) error {
	cPath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(
		C.krun_set_console_output(C.uint32_t(c.id), cPath),
		"krun_set_console_output",
	)
}

// DisableImplicitConsole prevents libkrun from creating an implicit console device.
// Any needed console devices must be added manually via other methods.
func (c *Context) DisableImplicitConsole() error {
	return checkRet(
		C.krun_disable_implicit_console(C.uint32_t(c.id)),
		"krun_disable_implicit_console",
	)
}

// SetKernelConsole sets the console= parameter in the kernel command line.
func (c *Context) SetKernelConsole(consoleID string) error {
	cID := C.CString(consoleID)
	defer C.free(unsafe.Pointer(cID))
	return checkRet(
		C.krun_set_kernel_console(C.uint32_t(c.id), cID),
		"krun_set_kernel_console",
	)
}

// AddVirtioConsoleDefault adds a virtio-console device with automatic detection.
// If the file descriptors are TTYs, a single console port is created.
// For non-TTY file descriptors, additional ports are created for stdin/stdout/stderr.
func (c *Context) AddVirtioConsoleDefault(cfg VirtioConsoleConfig) error {
	return checkRet(
		C.krun_add_virtio_console_default(
			C.uint32_t(c.id), C.int(cfg.InputFD), C.int(cfg.OutputFD), C.int(cfg.ErrFD),
		),
		"krun_add_virtio_console_default",
	)
}

// AddSerialConsoleDefault adds a legacy serial device.
func (c *Context) AddSerialConsoleDefault(cfg SerialConsoleConfig) error {
	return checkRet(
		C.krun_add_serial_console_default(C.uint32_t(c.id), C.int(cfg.InputFD), C.int(cfg.OutputFD)),
		"krun_add_serial_console_default",
	)
}

// AddVirtioConsoleMultiport creates a multi-port virtio-console device.
// Returns the console ID for use with [Context.AddConsolePortTTY] and
// [Context.AddConsolePortInOut].
func (c *Context) AddVirtioConsoleMultiport() (uint32, error) {
	ret := C.krun_add_virtio_console_multiport(C.uint32_t(c.id))
	if ret < 0 {
		return 0, retError(ret, "krun_add_virtio_console_multiport")
	}
	return uint32(ret), nil
}

// AddConsolePortTTY adds a TTY port to a multi-port virtio-console device.
// The port is marked with VIRTIO_CONSOLE_CONSOLE_PORT, enabling window resize support.
// Name identifies the port in the guest (can be "").
func (c *Context) AddConsolePortTTY(cfg ConsolePortTTYConfig) error {
	cName := C.CString(cfg.Name)
	defer C.free(unsafe.Pointer(cName))
	return checkRet(
		C.krun_add_console_port_tty(
			C.uint32_t(c.id), C.uint32_t(cfg.ConsoleID), cName, C.int(cfg.TTYFD),
		),
		"krun_add_console_port_tty",
	)
}

// AddConsolePortInOut adds a generic I/O port to a multi-port virtio-console device.
// The port does NOT support console features like window resize.
// Name identifies the port in the guest (can be "").
func (c *Context) AddConsolePortInOut(cfg ConsolePortInOutConfig) error {
	cName := C.CString(cfg.Name)
	defer C.free(unsafe.Pointer(cName))
	return checkRet(
		C.krun_add_console_port_inout(
			C.uint32_t(c.id), C.uint32_t(cfg.ConsoleID), cName,
			C.int(cfg.InputFD), C.int(cfg.OutputFD),
		),
		"krun_add_console_port_inout",
	)
}
