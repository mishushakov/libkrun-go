//go:build krun_tee

package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// SetTEEConfigFile sets the path to the TEE configuration file.
// Only available in libkrun-sev.
func (c *Context) SetTEEConfigFile(filepath string) error {
	cPath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(
		C.krun_set_tee_config_file(C.uint32_t(c.id), cPath),
		"krun_set_tee_config_file",
	)
}
