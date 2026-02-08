package krun

import (
	"os"
	"runtime"
	"testing"
)

func TestSetVMConfig(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.SetVMConfig(VMConfig{NumVCPUs: 2, RAMMiB: 512}); err != nil {
		t.Fatal(err)
	}
}

func TestSetRoot(t *testing.T) {
	ctx := newTestContext(t)
	dir := t.TempDir()
	if err := ctx.SetRoot(dir); err != nil {
		t.Fatal(err)
	}
}

func TestSetNestedVirt(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("SetNestedVirt(true) is only supported on macOS")
	}
	ctx := newTestContext(t)
	for _, enabled := range []bool{true, false} {
		if err := ctx.SetNestedVirt(enabled); err != nil {
			t.Errorf("SetNestedVirt(%v) = %v", enabled, err)
		}
	}
}

func TestSplitIRQChip(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("SplitIRQChip is only supported on Linux")
	}
	ctx := newTestContext(t)
	for _, enable := range []bool{true, false} {
		if err := ctx.SplitIRQChip(enable); err != nil {
			t.Errorf("SplitIRQChip(%v) = %v", enable, err)
		}
	}
}

func TestSetUID(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.SetUID(uint32(os.Getuid())); err != nil {
		t.Fatal(err)
	}
}

func TestSetGID(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.SetGID(uint32(os.Getgid())); err != nil {
		t.Fatal(err)
	}
}

func TestSetSMBIOSOEMStrings(t *testing.T) {
	ctx := newTestContext(t)

	t.Run("populated", func(t *testing.T) {
		err := ctx.SetSMBIOSOEMStrings([]string{"key1=val1", "key2=val2"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty", func(t *testing.T) {
		err := ctx.SetSMBIOSOEMStrings([]string{})
		if err != nil {
			t.Fatal(err)
		}
	})
}
