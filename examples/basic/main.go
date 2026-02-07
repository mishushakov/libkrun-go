// basic demonstrates the simplest libkrun usage: running a command
// inside a microVM backed by a host directory as root filesystem.
//
// Usage:
//
//	go run . /path/to/rootfs /bin/uname -a
package main

import (
	"fmt"
	"os"

	"e2b.dev/libkrun-go/krun"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <rootfs> <command> [args...]\n", os.Args[0])
		os.Exit(1)
	}

	rootfs := os.Args[1]
	execPath := os.Args[2]
	argv := os.Args[2:]

	if err := run(rootfs, execPath, argv); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(rootfs, execPath string, argv []string) error {
	// Enable info-level logging for visibility.
	if err := krun.SetLogLevel(krun.LogLevelInfo); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	// Create a new VM configuration context.
	ctx, err := krun.CreateContext()
	if err != nil {
		return fmt.Errorf("create context: %w", err)
	}

	// Configure 2 vCPUs and 512 MiB of RAM.
	if err := ctx.SetVMConfig(2, 512); err != nil {
		return fmt.Errorf("set vm config: %w", err)
	}

	// Use a host directory as the root filesystem.
	if err := ctx.SetRoot(rootfs); err != nil {
		return fmt.Errorf("set root: %w", err)
	}

	// Set the executable to run inside the VM.
	// Passing nil for envp auto-inherits the host environment.
	if err := ctx.SetExec(execPath, argv, nil); err != nil {
		return fmt.Errorf("set exec: %w", err)
	}

	// Start the VM. This call does not return on success — the process
	// exits with the guest workload's exit code.
	if err := ctx.StartEnter(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	return nil
}
