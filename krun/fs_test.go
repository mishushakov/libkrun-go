package krun

import "testing"

func TestAddVirtioFS(t *testing.T) {
	ctx := newTestContext(t)
	dir := t.TempDir()
	if err := ctx.AddVirtioFS("myfs", dir); err != nil {
		t.Fatal(err)
	}
}

func TestAddVirtioFS2(t *testing.T) {
	ctx := newTestContext(t)
	dir := t.TempDir()
	if err := ctx.AddVirtioFS2("myfs", dir, 256*1024*1024); err != nil {
		t.Fatal(err)
	}
}
