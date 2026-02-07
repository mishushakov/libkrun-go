//go:build !krun_blk

package krun

import "syscall"

// AddDisk adds a raw disk image as a partition for the microVM.
// Requires building with -tags krun_blk.
func (c *Context) AddDisk(blockID, diskPath string, readOnly bool) error {
	return &Error{Func: "krun_add_disk", Errno: syscall.ENOSYS}
}

// AddDisk2 adds a disk image with an explicit format as a partition for the microVM.
// Requires building with -tags krun_blk.
func (c *Context) AddDisk2(blockID, diskPath string, format DiskFormat, readOnly bool) error {
	return &Error{Func: "krun_add_disk2", Errno: syscall.ENOSYS}
}

// AddDisk3 adds a disk image with format, direct I/O, and sync mode options.
// Requires building with -tags krun_blk.
func (c *Context) AddDisk3(blockID, diskPath string, format DiskFormat, readOnly, directIO bool, syncMode SyncMode) error {
	return &Error{Func: "krun_add_disk3", Errno: syscall.ENOSYS}
}

// SetRootDiskRemount configures a block device as the root filesystem.
// Requires building with -tags krun_blk.
func (c *Context) SetRootDiskRemount(device, fstype, options string) error {
	return &Error{Func: "krun_set_root_disk_remount", Errno: syscall.ENOSYS}
}
