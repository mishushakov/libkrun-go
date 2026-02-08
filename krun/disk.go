//go:build krun_blk

package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// AddDisk adds a disk image as a partition for the microVM.
func (c *Context) AddDisk(cfg DiskConfig) error {
	cBlockID := C.CString(cfg.BlockID)
	defer C.free(unsafe.Pointer(cBlockID))
	cDiskPath := C.CString(cfg.Path)
	defer C.free(unsafe.Pointer(cDiskPath))
	return checkRet(
		C.krun_add_disk3(
			C.uint32_t(c.id), cBlockID, cDiskPath,
			C.uint32_t(cfg.Format), C.bool(cfg.ReadOnly), C.bool(cfg.DirectIO), C.uint32_t(cfg.SyncMode),
		),
		"krun_add_disk3",
	)
}

// SetRootDiskRemount configures a block device as the root filesystem.
// Device must refer to a previously configured block device (e.g., "/dev/vda1").
// FSType is the filesystem type (e.g., "ext4" or "auto"). Pass "" for NULL.
// Options is a comma-separated list of mount options. Pass "" for NULL.
func (c *Context) SetRootDiskRemount(cfg RootDiskRemountConfig) error {
	cDevice := C.CString(cfg.Device)
	defer C.free(unsafe.Pointer(cDevice))

	var cFstype *C.char
	if cfg.FSType != "" {
		cFstype = C.CString(cfg.FSType)
		defer C.free(unsafe.Pointer(cFstype))
	}

	var cOptions *C.char
	if cfg.Options != "" {
		cOptions = C.CString(cfg.Options)
		defer C.free(unsafe.Pointer(cOptions))
	}

	return checkRet(
		C.krun_set_root_disk_remount(C.uint32_t(c.id), cDevice, cFstype, cOptions),
		"krun_set_root_disk_remount",
	)
}
