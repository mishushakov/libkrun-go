// Package krun provides Go bindings for libkrun, a dynamic library for
// creating lightweight microVMs using KVM (Linux) or HVF (macOS/ARM64).
//
// The typical workflow is:
//  1. Create a context with [CreateContext]
//  2. Configure the VM using methods on [Context]
//  3. Call [Context.StartEnter] to launch the microVM
package krun

/*
#cgo CFLAGS: -I${SRCDIR}/../libkrun/include
#cgo LDFLAGS: -lkrun
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"
)

// Context represents a libkrun VM configuration context.
type Context struct {
	id uint32
}

// ID returns the underlying context ID.
func (c *Context) ID() uint32 {
	return c.id
}

// Error represents an error returned by the libkrun library.
type Error struct {
	// Func is the libkrun C function that failed.
	Func string
	// Errno is the system error number returned by the function.
	Errno syscall.Errno
}

func (e *Error) Error() string {
	return fmt.Sprintf("krun: %s: %s", e.Func, e.Errno)
}

func (e *Error) Unwrap() error {
	return e.Errno
}

func checkRet(ret C.int32_t, fn string) error {
	if ret >= 0 {
		return nil
	}
	return &Error{Func: fn, Errno: syscall.Errno(-ret)}
}

func retError(ret C.int32_t, fn string) error {
	return &Error{Func: fn, Errno: syscall.Errno(-ret)}
}

// SetLogLevel sets the log level for the library.
func SetLogLevel(level LogLevel) error {
	return checkRet(C.krun_set_log_level(C.uint32_t(level)), "krun_set_log_level")
}

// InitLog initializes logging for the library.
// Use [LogTargetDefault] as targetFD for the default log target (stderr).
func InitLog(targetFD int, level LogLevel, style LogStyle, options uint32) error {
	return checkRet(
		C.krun_init_log(C.int(targetFD), C.uint32_t(level), C.uint32_t(style), C.uint32_t(options)),
		"krun_init_log",
	)
}

// CreateContext creates a new VM configuration context.
func CreateContext() (*Context, error) {
	ret := C.krun_create_ctx()
	if ret < 0 {
		return nil, retError(ret, "krun_create_ctx")
	}
	return &Context{id: uint32(ret)}, nil
}

// Free releases the configuration context.
func (c *Context) Free() error {
	return checkRet(C.krun_free_ctx(C.uint32_t(c.id)), "krun_free_ctx")
}

// StartEnter starts and enters the microVM. This function consumes the context.
// It only returns if an error occurs before starting the microVM. Otherwise,
// the VMM calls exit() with the workload's exit code once the VM shuts down.
func (c *Context) StartEnter() error {
	return checkRet(C.krun_start_enter(C.uint32_t(c.id)), "krun_start_enter")
}

// HasFeature checks if a specific feature was enabled at build time.
func HasFeature(feature Feature) (bool, error) {
	ret := C.krun_has_feature(C.uint64_t(feature))
	if ret < 0 {
		return false, retError(ret, "krun_has_feature")
	}
	return ret == 1, nil
}

// GetMaxVCPUs returns the maximum number of vCPUs supported by the hypervisor.
func GetMaxVCPUs() (int, error) {
	ret := C.krun_get_max_vcpus()
	if ret < 0 {
		return 0, retError(ret, "krun_get_max_vcpus")
	}
	return int(ret), nil
}

// CheckNestedVirt checks if nested virtualization is supported (macOS only).
func CheckNestedVirt() (bool, error) {
	ret := C.krun_check_nested_virt()
	if ret < 0 {
		return false, retError(ret, "krun_check_nested_virt")
	}
	return ret == 1, nil
}

// stringsToCArray converts a Go string slice to a C null-terminated string array.
// Returns nil for nil input (which maps to C NULL).
// Returns a null-terminated empty array for empty non-nil input.
// The caller must call freeCStringArray to release the memory.
func stringsToCArray(strs []string) **C.char {
	if strs == nil {
		return nil
	}
	n := len(strs)
	ptrSize := unsafe.Sizeof((*C.char)(nil))
	arr := (**C.char)(C.malloc(C.size_t(uintptr(n+1) * ptrSize)))
	slice := unsafe.Slice(arr, n+1)
	for i, s := range strs {
		slice[i] = C.CString(s)
	}
	slice[n] = nil
	return arr
}

func freeCStringArray(arr **C.char, n int) {
	if arr == nil {
		return
	}
	slice := unsafe.Slice(arr, n+1)
	for i := range n {
		C.free(unsafe.Pointer(slice[i]))
	}
	C.free(unsafe.Pointer(arr))
}
