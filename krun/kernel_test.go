package krun

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetFirmware(t *testing.T) {
	ctx := newTestContext(t)
	// Create a temp file to act as firmware
	fwPath := filepath.Join(t.TempDir(), "firmware.bin")
	if err := os.WriteFile(fwPath, []byte("fake firmware"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ctx.SetFirmware(fwPath); err != nil {
		t.Fatal(err)
	}
}

func TestSetKernel(t *testing.T) {
	ctx := newTestContext(t)
	kernelPath := filepath.Join(t.TempDir(), "kernel.bin")
	if err := os.WriteFile(kernelPath, []byte("fake kernel"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("minimal", func(t *testing.T) {
		err := ctx.SetKernel(KernelConfig{Path: kernelPath})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with_cmdline", func(t *testing.T) {
		err := ctx.SetKernel(KernelConfig{Path: kernelPath, Cmdline: "console=ttyS0"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with_initramfs", func(t *testing.T) {
		initrdPath := filepath.Join(t.TempDir(), "initramfs.img")
		if err := os.WriteFile(initrdPath, []byte("fake initramfs"), 0644); err != nil {
			t.Fatal(err)
		}
		err := ctx.SetKernel(KernelConfig{
			Path:      kernelPath,
			Format:    KernelFormatELF,
			Initramfs: initrdPath,
			Cmdline:   "root=/dev/vda",
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}
