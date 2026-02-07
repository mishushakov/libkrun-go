//go:build krun_blk

package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// AddDisk adds a raw disk image as a partition for the microVM.
func (c *Context) AddDisk(blockID, diskPath string, readOnly bool) error {
	cBlockID := C.CString(blockID)
	defer C.free(unsafe.Pointer(cBlockID))
	cDiskPath := C.CString(diskPath)
	defer C.free(unsafe.Pointer(cDiskPath))
	return checkRet(
		C.krun_add_disk(C.uint32_t(c.id), cBlockID, cDiskPath, C.bool(readOnly)),
		"krun_add_disk",
	)
}

// AddDisk2 adds a disk image with an explicit format as a partition for the microVM.
//
// Security note: Non-raw images can reference other files. Only use non-raw formats
// with fully trusted images. See the libkrun documentation for details.
func (c *Context) AddDisk2(blockID, diskPath string, format DiskFormat, readOnly bool) error {
	cBlockID := C.CString(blockID)
	defer C.free(unsafe.Pointer(cBlockID))
	cDiskPath := C.CString(diskPath)
	defer C.free(unsafe.Pointer(cDiskPath))
	return checkRet(
		C.krun_add_disk2(C.uint32_t(c.id), cBlockID, cDiskPath, C.uint32_t(format), C.bool(readOnly)),
		"krun_add_disk2",
	)
}

// AddDisk3 adds a disk image with format, direct I/O, and sync mode options.
//
// Security note: Non-raw images can reference other files. Only use non-raw formats
// with fully trusted images. See the libkrun documentation for details.
func (c *Context) AddDisk3(blockID, diskPath string, format DiskFormat, readOnly, directIO bool, syncMode SyncMode) error {
	cBlockID := C.CString(blockID)
	defer C.free(unsafe.Pointer(cBlockID))
	cDiskPath := C.CString(diskPath)
	defer C.free(unsafe.Pointer(cDiskPath))
	return checkRet(
		C.krun_add_disk3(
			C.uint32_t(c.id), cBlockID, cDiskPath,
			C.uint32_t(format), C.bool(readOnly), C.bool(directIO), C.uint32_t(syncMode),
		),
		"krun_add_disk3",
	)
}

// SetRootDiskRemount configures a block device as the root filesystem.
// device must refer to a previously configured block device (e.g., "/dev/vda1").
// fstype is the filesystem type (e.g., "ext4" or "auto"). Pass "" for NULL.
// options is a comma-separated list of mount options. Pass "" for NULL.
func (c *Context) SetRootDiskRemount(device, fstype, options string) error {
	cDevice := C.CString(device)
	defer C.free(unsafe.Pointer(cDevice))

	var cFstype *C.char
	if fstype != "" {
		cFstype = C.CString(fstype)
		defer C.free(unsafe.Pointer(cFstype))
	}

	var cOptions *C.char
	if options != "" {
		cOptions = C.CString(options)
		defer C.free(unsafe.Pointer(cOptions))
	}

	return checkRet(
		C.krun_set_root_disk_remount(C.uint32_t(c.id), cDevice, cFstype, cOptions),
		"krun_set_root_disk_remount",
	)
}
