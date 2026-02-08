package krun

import (
	"errors"
	"syscall"
	"testing"
)

// Stub tests — verify behavior when built without -tags krun_blk.

func TestAddDisk_Stub(t *testing.T) {
	ctx := newTestContext(t)
	err := ctx.AddDisk("vda", "/tmp/disk.img", false)
	if err == nil {
		return // built with krun_blk tag
	}
	if !errors.Is(err, syscall.ENOSYS) {
		t.Fatalf("expected ENOSYS, got %v", err)
	}
}

func TestAddDisk2_Stub(t *testing.T) {
	ctx := newTestContext(t)
	err := ctx.AddDisk2("vda", "/tmp/disk.img", DiskFormatQcow2, false)
	if err == nil {
		return
	}
	if !errors.Is(err, syscall.ENOSYS) {
		t.Fatalf("expected ENOSYS, got %v", err)
	}
}

func TestAddDisk3_Stub(t *testing.T) {
	ctx := newTestContext(t)
	err := ctx.AddDisk3("vda", "/tmp/disk.img", DiskFormatRaw, true, false, SyncRelaxed)
	if err == nil {
		return
	}
	if !errors.Is(err, syscall.ENOSYS) {
		t.Fatalf("expected ENOSYS, got %v", err)
	}
}

func TestSetRootDiskRemount_Stub(t *testing.T) {
	ctx := newTestContext(t)
	err := ctx.SetRootDiskRemount("/dev/vda", "ext4", "rw")
	if err == nil {
		return
	}
	if !errors.Is(err, syscall.ENOSYS) {
		t.Fatalf("expected ENOSYS, got %v", err)
	}
}
