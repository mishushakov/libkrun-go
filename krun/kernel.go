package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// KernelConfig configures the kernel to be loaded in the microVM.
type KernelConfig struct {
	Path      string
	Format    KernelFormat // 0 = KernelFormatRaw
	Initramfs string       // "" = none
	Cmdline   string       // "" = none
}

// SetFirmware sets the path to the firmware to be loaded into the microVM.
func (c *Context) SetFirmware(firmwarePath string) error {
	cPath := C.CString(firmwarePath)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(C.krun_set_firmware(C.uint32_t(c.id), cPath), "krun_set_firmware")
}

// SetKernel configures the kernel to be loaded in the microVM.
func (c *Context) SetKernel(cfg KernelConfig) error {
	cKernel := C.CString(cfg.Path)
	defer C.free(unsafe.Pointer(cKernel))

	var cInitramfs *C.char
	if cfg.Initramfs != "" {
		cInitramfs = C.CString(cfg.Initramfs)
		defer C.free(unsafe.Pointer(cInitramfs))
	}

	var cCmdline *C.char
	if cfg.Cmdline != "" {
		cCmdline = C.CString(cfg.Cmdline)
		defer C.free(unsafe.Pointer(cCmdline))
	}

	return checkRet(
		C.krun_set_kernel(
			C.uint32_t(c.id), cKernel, C.uint32_t(cfg.Format), cInitramfs, cCmdline,
		),
		"krun_set_kernel",
	)
}
