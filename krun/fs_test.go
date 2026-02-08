package krun

import "testing"

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
