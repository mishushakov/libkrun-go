package krun

import (
	"path/filepath"
	"testing"
)

func TestSetConsoleOutput(t *testing.T) {
	ctx := newTestContext(t)
	outPath := filepath.Join(t.TempDir(), "console.log")
	if err := ctx.SetConsoleOutput(outPath); err != nil {
		t.Fatal(err)
	}
}

func TestDisableImplicitConsole(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.DisableImplicitConsole(); err != nil {
		t.Fatal(err)
	}
}

func TestSetKernelConsole(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.SetKernelConsole("hvc0"); err != nil {
		t.Fatal(err)
	}
}

func TestAddVirtioConsoleMultiport(t *testing.T) {
	ctx := newTestContext(t)
	id, err := ctx.AddVirtioConsoleMultiport()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("console ID: %d", id)
}
