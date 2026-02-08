package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// VirtioFSConfig configures a virtio-fs device.
type VirtioFSConfig struct {
	Tag     string
	Path    string
	ShmSize uint64 // 0 = libkrun default
}

// AddVirtioFS adds a virtio-fs device pointing to a host directory.
func (c *Context) AddVirtioFS(cfg VirtioFSConfig) error {
	cTag := C.CString(cfg.Tag)
	defer C.free(unsafe.Pointer(cTag))
	cPath := C.CString(cfg.Path)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(
		C.krun_add_virtiofs2(C.uint32_t(c.id), cTag, cPath, C.uint64_t(cfg.ShmSize)),
		"krun_add_virtiofs2",
	)
}
