// vm-with-disk demonstrates booting a microVM from a disk image with a
// custom kernel, virtio-fs shared directory, and console setup.
//
// Usage:
//
//	go run . -kernel /path/to/vmlinux -disk /path/to/rootfs.ext4
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mishushakov/libkrun-go/krun"
)

func main() {
	var (
		kernel  = flag.String("kernel", "", "path to kernel image")
		disk    = flag.String("disk", "", "path to root disk image (ext4)")
		format  = flag.String("format", "raw", "disk format: raw, qcow2, vmdk")
		shared  = flag.String("shared", "", "host directory to share via virtio-fs")
		vcpus   = flag.Int("vcpus", 2, "number of vCPUs")
		ram     = flag.Int("ram", 1024, "RAM in MiB")
		cmdline = flag.String("cmdline", "console=hvc0 root=/dev/vda1 rw", "kernel command line")
	)
	flag.Parse()

	if *kernel == "" || *disk == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*kernel, *disk, *format, *shared, *vcpus, *ram, *cmdline); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func parseDiskFormat(s string) (krun.DiskFormat, error) {
	switch s {
	case "raw":
		return krun.DiskFormatRaw, nil
	case "qcow2":
		return krun.DiskFormatQcow2, nil
	case "vmdk":
		return krun.DiskFormatVmdk, nil
	default:
		return 0, fmt.Errorf("unknown disk format: %q", s)
	}
}

func run(kernel, disk, format, shared string, vcpus, ram int, cmdline string) error {
	if err := krun.SetLogLevel(krun.LogLevelInfo); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	ctx, err := krun.CreateContext()
	if err != nil {
		return fmt.Errorf("create context: %w", err)
	}

	if err := ctx.SetVMConfig(krun.VMConfig{NumVCPUs: uint8(vcpus), RAMMiB: uint32(ram)}); err != nil {
		return fmt.Errorf("set vm config: %w", err)
	}

	// Load the kernel.
	if err := ctx.SetKernel(krun.KernelConfig{Path: kernel, Cmdline: cmdline}); err != nil {
		return fmt.Errorf("set kernel: %w", err)
	}

	// Attach the root disk.
	diskFmt, err := parseDiskFormat(format)
	if err != nil {
		return err
	}
	if err := ctx.AddDisk(krun.DiskConfig{BlockID: "vda", Path: disk, Format: diskFmt}); err != nil {
		return fmt.Errorf("add disk: %w", err)
	}

	// Remount the disk as root filesystem.
	if err := ctx.SetRootDiskRemount(krun.RootDiskRemountConfig{Device: "/dev/vda1", FSType: "ext4"}); err != nil {
		return fmt.Errorf("set root disk remount: %w", err)
	}

	// Optionally share a host directory into the guest.
	if shared != "" {
		if err := ctx.AddVirtioFS(krun.VirtioFSConfig{Tag: "shared", Path: shared}); err != nil {
			return fmt.Errorf("add virtiofs: %w", err)
		}
	}

	// Set up a console using stdin/stdout/stderr.
	if err := ctx.AddVirtioConsoleDefault(krun.VirtioConsoleConfig{
		InputFD:  int(os.Stdin.Fd()),
		OutputFD: int(os.Stdout.Fd()),
		ErrFD:    int(os.Stderr.Fd()),
	}); err != nil {
		return fmt.Errorf("add console: %w", err)
	}

	if err := ctx.StartEnter(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	return nil
}
