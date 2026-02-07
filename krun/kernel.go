package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// SetFirmware sets the path to the firmware to be loaded into the microVM.
func (c *Context) SetFirmware(firmwarePath string) error {
	cPath := C.CString(firmwarePath)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(C.krun_set_firmware(C.uint32_t(c.id), cPath), "krun_set_firmware")
}

// SetKernel configures the kernel to be loaded in the microVM.
// kernelPath is the path to the kernel image.
// format specifies the kernel image format.
// initramfs is the path to the initramfs (pass "" for none).
// cmdline is the kernel command line (pass "" for none).
func (c *Context) SetKernel(kernelPath string, format KernelFormat, initramfs, cmdline string) error {
	cKernel := C.CString(kernelPath)
	defer C.free(unsafe.Pointer(cKernel))

	var cInitramfs *C.char
	if initramfs != "" {
		cInitramfs = C.CString(initramfs)
		defer C.free(unsafe.Pointer(cInitramfs))
	}

	var cCmdline *C.char
	if cmdline != "" {
		cCmdline = C.CString(cmdline)
		defer C.free(unsafe.Pointer(cCmdline))
	}

	return checkRet(
		C.krun_set_kernel(
			C.uint32_t(c.id), cKernel, C.uint32_t(format), cInitramfs, cCmdline,
		),
		"krun_set_kernel",
	)
}
