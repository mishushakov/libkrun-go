//go:build !krun_blk

package krun

import "syscall"

// AddDisk adds a disk image as a partition for the microVM.
// Requires building with -tags krun_blk.
func (c *Context) AddDisk(cfg DiskConfig) error {
	return &Error{Func: "krun_add_disk3", Errno: syscall.ENOSYS}
}

// SetRootDiskRemount configures a block device as the root filesystem.
// Requires building with -tags krun_blk.
func (c *Context) SetRootDiskRemount(cfg RootDiskRemountConfig) error {
	return &Error{Func: "krun_set_root_disk_remount", Errno: syscall.ENOSYS}
}
