package krun

import (
	"errors"
	"syscall"
	"testing"
)

// Stub tests — verify behavior when built without -tags krun_blk.

func TestAddDisk_Stub(t *testing.T) {
	ctx := newTestContext(t)
	err := ctx.AddDisk(DiskConfig{
		BlockID: "vda",
		Path:    "/tmp/disk.img",
	})
	if err == nil {
		return // built with krun_blk tag
	}
	if !errors.Is(err, syscall.ENOSYS) {
		t.Fatalf("expected ENOSYS, got %v", err)
	}
}

func TestAddDisk_WithOptions(t *testing.T) {
	ctx := newTestContext(t)
	err := ctx.AddDisk(DiskConfig{
		BlockID:  "vda",
		Path:     "/tmp/disk.img",
		Format:   DiskFormatQcow2,
		ReadOnly: true,
		DirectIO: false,
		SyncMode: SyncRelaxed,
	})
	if err == nil {
		return
	}
	if !errors.Is(err, syscall.ENOSYS) {
		t.Fatalf("expected ENOSYS, got %v", err)
	}
}

func TestSetRootDiskRemount_Stub(t *testing.T) {
	ctx := newTestContext(t)
	err := ctx.SetRootDiskRemount(RootDiskRemountConfig{Device: "/dev/vda", FSType: "ext4", Options: "rw"})
	if err == nil {
		return
	}
	// ENOSYS when built without krun_blk; EINVAL when built with krun_blk
	// but no disk has been configured on the context.
	if !errors.Is(err, syscall.ENOSYS) && !errors.Is(err, syscall.EINVAL) {
		t.Fatalf("expected ENOSYS or EINVAL, got %v", err)
	}
}
