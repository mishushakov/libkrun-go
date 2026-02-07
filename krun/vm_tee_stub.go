//go:build !krun_tee

package krun

import "syscall"

// SetTEEConfigFile sets the path to the TEE configuration file.
// Requires building with -tags krun_tee.
func (c *Context) SetTEEConfigFile(filepath string) error {
	return &Error{Func: "krun_set_tee_config_file", Errno: syscall.ENOSYS}
}
