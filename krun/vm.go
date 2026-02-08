package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// SetVMConfig sets the basic configuration parameters for the microVM.
func (c *Context) SetVMConfig(cfg VMConfig) error {
	return checkRet(
		C.krun_set_vm_config(C.uint32_t(c.id), C.uint8_t(cfg.NumVCPUs), C.uint32_t(cfg.RAMMiB)),
		"krun_set_vm_config",
	)
}

// SetRoot sets the path to be used as root for the microVM.
// Not available in libkrun-SEV.
func (c *Context) SetRoot(rootPath string) error {
	cPath := C.CString(rootPath)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(C.krun_set_root(C.uint32_t(c.id), cPath), "krun_set_root")
}

// SetNestedVirt enables or disables nested virtualization (macOS only).
func (c *Context) SetNestedVirt(enabled bool) error {
	return checkRet(
		C.krun_set_nested_virt(C.uint32_t(c.id), C.bool(enabled)),
		"krun_set_nested_virt",
	)
}

// SplitIRQChip specifies whether to split IRQCHIP responsibilities
// between the host and the guest.
func (c *Context) SplitIRQChip(enable bool) error {
	return checkRet(
		C.krun_split_irqchip(C.uint32_t(c.id), C.bool(enable)),
		"krun_split_irqchip",
	)
}

// SetUID sets the user ID before the microVM is started.
// Useful when root privileges are needed to open devices but the VM
// should not run as root.
func (c *Context) SetUID(uid uint32) error {
	return checkRet(C.krun_setuid(C.uint32_t(c.id), C.uid_t(uid)), "krun_setuid")
}

// SetGID sets the group ID before the microVM is started.
func (c *Context) SetGID(gid uint32) error {
	return checkRet(C.krun_setgid(C.uint32_t(c.id), C.gid_t(gid)), "krun_setgid")
}

// SetSMBIOSOEMStrings sets the SMBIOS OEM Strings.
func (c *Context) SetSMBIOSOEMStrings(oemStrings []string) error {
	cArr := stringsToCArray(oemStrings)
	defer freeCStringArray(cArr, len(oemStrings))
	return checkRet(
		C.krun_set_smbios_oem_strings(C.uint32_t(c.id), cArr),
		"krun_set_smbios_oem_strings",
	)
}

// GetShutdownEventFD returns a file descriptor that can be used to signal
// the guest to shut down. Only available in libkrun-efi.
// Must be called before [Context.StartEnter].
func (c *Context) GetShutdownEventFD() (int, error) {
	ret := C.krun_get_shutdown_eventfd(C.uint32_t(c.id))
	if ret < 0 {
		return 0, retError(ret, "krun_get_shutdown_eventfd")
	}
	return int(ret), nil
}
