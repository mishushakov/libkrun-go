package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// SetWorkdir sets the working directory for the executable to be run
// inside the microVM. The path is relative to the root configured with [Context.SetRoot].
func (c *Context) SetWorkdir(workdirPath string) error {
	cPath := C.CString(workdirPath)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(C.krun_set_workdir(C.uint32_t(c.id), cPath), "krun_set_workdir")
}

// SetExec sets the executable path, arguments, and environment variables
// for the process to run inside the microVM.
//
// Path is relative to the root configured with [Context.SetRoot].
// Args is the argument list (Args[0] is typically the program name).
// Env is the environment variables (e.g., "KEY=value"). Pass nil to
// auto-generate from the current process environment.
func (c *Context) SetExec(cfg ExecConfig) error {
	cExec := C.CString(cfg.Path)
	defer C.free(unsafe.Pointer(cExec))

	cArgv := stringsToCArray(cfg.Args)
	defer freeCStringArray(cArgv, len(cfg.Args))

	cEnvp := stringsToCArray(cfg.Env)
	if cfg.Env != nil {
		defer freeCStringArray(cEnvp, len(cfg.Env))
	}

	return checkRet(
		C.krun_set_exec(C.uint32_t(c.id), cExec, cArgv, cEnvp),
		"krun_set_exec",
	)
}

// SetEnv sets environment variables for the executable.
// Pass nil to auto-generate from the current process environment.
func (c *Context) SetEnv(envp []string) error {
	cEnvp := stringsToCArray(envp)
	if envp != nil {
		defer freeCStringArray(cEnvp, len(envp))
	}
	return checkRet(C.krun_set_env(C.uint32_t(c.id), cEnvp), "krun_set_env")
}

// SetRlimits configures resource limits to be set in the guest before
// starting the executable. Each entry has the format "RESOURCE=RLIM_CUR:RLIM_MAX".
func (c *Context) SetRlimits(rlimits []string) error {
	cArr := stringsToCArray(rlimits)
	defer freeCStringArray(cArr, len(rlimits))
	return checkRet(C.krun_set_rlimits(C.uint32_t(c.id), cArr), "krun_set_rlimits")
}
