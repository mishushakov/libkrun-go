package krun

import (
	"errors"
	"syscall"
	"testing"
)

func TestAddVirtioFS(t *testing.T) {
	ctx := newTestContext(t)
	dir := t.TempDir()
	if err := ctx.AddVirtioFS(VirtioFSConfig{Tag: "myfs", Path: dir}); err != nil {
		t.Fatal(err)
	}
}

func TestAddVirtioFS_WithShmSize(t *testing.T) {
	ctx := newTestContext(t)
	dir := t.TempDir()
	if err := ctx.AddVirtioFS(VirtioFSConfig{Tag: "myfs", Path: dir, ShmSize: 256 * 1024 * 1024}); err != nil {
		t.Fatal(err)
	}
}

func TestAddVirtioFS_ReadOnly(t *testing.T) {
	ctx := newTestContext(t)
	dir := t.TempDir()
	if err := ctx.AddVirtioFS(VirtioFSConfig{Tag: "myfs", Path: dir, ReadOnly: true}); err != nil {
		t.Fatal(err)
	}
}

func TestAddVirtioFS_SimplifiedSemantics(t *testing.T) {
	ctx := newTestContext(t)
	dir := t.TempDir()
	if err := ctx.AddVirtioFS(VirtioFSConfig{Tag: "myfs", Path: dir, Semantics: FSSemanticsLinuxSimplified}); err != nil {
		t.Fatal(err)
	}
}

func TestAddVirtioFS_InvalidSemantics(t *testing.T) {
	ctx := newTestContext(t)
	dir := t.TempDir()
	err := ctx.AddVirtioFS(VirtioFSConfig{Tag: "myfs", Path: dir, Semantics: 99})
	if !errors.Is(err, syscall.EINVAL) {
		t.Fatalf("AddVirtioFS(Semantics: 99) = %v, want EINVAL", err)
	}
}

func TestAddVirtioFS_RootTag(t *testing.T) {
	ctx := newTestContext(t)
	dir := t.TempDir()
	if err := ctx.AddVirtioFS(VirtioFSConfig{Tag: FSRootTag, Path: dir, ReadOnly: true}); err != nil {
		t.Fatal(err)
	}
}

func TestFSAddOverlay(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.SetRoot(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	if err := ctx.FSAddOverlayDir(FSRootTag, "etc", 0o40755); err != nil {
		t.Fatal(err)
	}
	if err := ctx.FSAddOverlayFile(FSRootTag, "etc/hostname", []byte("krun\n"), 0o100644, false); err != nil {
		t.Fatal(err)
	}
	if err := ctx.FSAddOverlayFile(FSRootTag, "etc/empty", nil, 0o100644, true); err != nil {
		t.Fatal(err)
	}
}

func TestFSAddOverlay_Errors(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.SetRoot(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	if err := ctx.FSAddOverlayDir("nosuchtag", "dir", 0o40755); !errors.Is(err, syscall.ENOENT) {
		t.Errorf("FSAddOverlayDir(unknown tag) = %v, want ENOENT", err)
	}
	if err := ctx.FSAddOverlayFile(FSRootTag, "missing/file", []byte("x"), 0o100644, false); !errors.Is(err, syscall.ENOENT) {
		t.Errorf("FSAddOverlayFile(missing parent) = %v, want ENOENT", err)
	}
	if err := ctx.FSAddOverlayFile(FSRootTag, "file", []byte("x"), 0o100644, false); err != nil {
		t.Fatal(err)
	}
	if err := ctx.FSAddOverlayFile(FSRootTag, "file/child", []byte("x"), 0o100644, false); !errors.Is(err, syscall.ENOTDIR) {
		t.Errorf("FSAddOverlayFile(file parent) = %v, want ENOTDIR", err)
	}
}

func TestGetDefaultInit(t *testing.T) {
	hasInit, err := HasFeature(FeatureInitBlob)
	if err != nil {
		t.Fatal(err)
	}
	data, err := GetDefaultInit()
	if !hasInit {
		if !errors.Is(err, syscall.ENOTSUP) {
			t.Fatalf("GetDefaultInit() = %v, want ENOTSUP", err)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("GetDefaultInit() returned empty data")
	}
}
