package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// AddVirtioFS adds a virtio-fs device pointing to a host directory.
// tag identifies the filesystem in the guest.
// path is the full path to the host directory to expose.
func (c *Context) AddVirtioFS(tag, path string) error {
	cTag := C.CString(tag)
	defer C.free(unsafe.Pointer(cTag))
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(
		C.krun_add_virtiofs(C.uint32_t(c.id), cTag, cPath),
		"krun_add_virtiofs",
	)
}

// AddVirtioFS2 adds a virtio-fs device with a custom DAX window size.
// tag identifies the filesystem in the guest.
// path is the full path to the host directory to expose.
// shmSize is the DAX SHM window size in bytes.
func (c *Context) AddVirtioFS2(tag, path string, shmSize uint64) error {
	cTag := C.CString(tag)
	defer C.free(unsafe.Pointer(cTag))
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(
		C.krun_add_virtiofs2(C.uint32_t(c.id), cTag, cPath, C.uint64_t(shmSize)),
		"krun_add_virtiofs2",
	)
}
