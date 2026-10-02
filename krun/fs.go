package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// VirtioFSConfig configures a virtio-fs device.
type VirtioFSConfig struct {
	Tag       string
	Path      string
	ShmSize   uint64      // 0 = libkrun default
	ReadOnly  bool        // expose the filesystem as read-only to the guest
	Semantics FSSemantics // 0 = FSSemanticsLinuxComplete
}

// AddVirtioFS adds a virtio-fs device pointing to a host directory.
// Use [FSRootTag] as the tag to configure the root filesystem.
func (c *Context) AddVirtioFS(cfg VirtioFSConfig) error {
	cTag := C.CString(cfg.Tag)
	defer C.free(unsafe.Pointer(cTag))
	cPath := C.CString(cfg.Path)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(
		C.krun_add_virtiofs4(
			C.uint32_t(c.id), cTag, cPath, C.uint64_t(cfg.ShmSize),
			C.bool(cfg.ReadOnly), C.uint32_t(cfg.Semantics),
		),
		"krun_add_virtiofs4",
	)
}

// FSAddOverlayFile adds a virtual file, backed entirely by memory, to the
// virtio-fs device identified by fsTag (e.g. [FSRootTag]).
//
// path may contain '/' to place the file inside a virtual directory previously
// created with [Context.FSAddOverlayDir] (e.g. "etc/hostname"); all
// intermediate directories must already exist. mode holds the file mode bits
// (e.g. 0o100644 for a regular file). If oneShot is true, the file can only be
// looked up once.
//
// data is copied into memory that stays valid for the VM lifetime and is
// released by [Context.Free].
func (c *Context) FSAddOverlayFile(fsTag, path string, data []byte, mode uint32, oneShot bool) error {
	cTag := C.CString(fsTag)
	defer C.free(unsafe.Pointer(cTag))
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var cData unsafe.Pointer
	if len(data) > 0 {
		cData = C.CBytes(data)
	}
	ret := C.krun_fs_add_overlay_file(
		C.uint32_t(c.id), cTag, cPath, (*C.uint8_t)(cData), C.size_t(len(data)),
		C.uint32_t(mode), C.bool(oneShot),
	)
	if ret < 0 {
		C.free(cData)
		return retError(ret, "krun_fs_add_overlay_file")
	}
	if cData != nil {
		c.allocs = append(c.allocs, cData)
	}
	return nil
}

// FSAddOverlayDir adds an empty, read-only virtual directory to the virtio-fs
// device identified by fsTag (e.g. [FSRootTag]), useful as a mount point.
//
// path may contain '/' to nest inside an existing virtual directory
// (e.g. "usr/lib"); all intermediate directories must already exist.
// mode holds the directory mode bits (e.g. 0o40755).
func (c *Context) FSAddOverlayDir(fsTag, path string, mode uint32) error {
	cTag := C.CString(fsTag)
	defer C.free(unsafe.Pointer(cTag))
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	return checkRet(
		C.krun_fs_add_overlay_dir(C.uint32_t(c.id), cTag, cPath, C.uint32_t(mode)),
		"krun_fs_add_overlay_dir",
	)
}

// GetDefaultInit returns a copy of the built-in default init binary, the same
// one libkrun injects as /init.krun unless [Context.DisableImplicitInit] is
// called. Returns ENOTSUP if libkrun was built without the init-blob feature
// (see [FeatureInitBlob]).
func GetDefaultInit() ([]byte, error) {
	var data *C.uint8_t
	var n C.size_t
	if ret := C.krun_get_default_init(&data, &n); ret < 0 {
		return nil, retError(ret, "krun_get_default_init")
	}
	return C.GoBytes(unsafe.Pointer(data), C.int(n)), nil
}
