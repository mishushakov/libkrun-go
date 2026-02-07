# libkrun-go

Go bindings for [libkrun](https://github.com/containers/libkrun), a dynamic library for creating lightweight microVMs using KVM (Linux) or HVF (macOS/ARM64).

## Installation

```bash
go get e2b.dev/libkrun-go/krun
```

libkrun must be installed on your system. The bundled header in `libkrun/include/libkrun.h` is used at build time, but the shared library (`libkrun.so` or `libkrun.dylib`) must be available to the linker.

### Building libkrun from source

```bash
cd libkrun
make
sudo make install
```

### macOS

libkrunfw (the firmware) is available via Homebrew:

```bash
brew install libkrunfw
```

libkrun itself must be built from source.

## Usage

```go
package main

import (
	"fmt"
	"os"

	"e2b.dev/libkrun-go/krun"
)

func main() {
	// Create a new VM configuration context.
	ctx, err := krun.CreateContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Configure 2 vCPUs and 512 MiB of RAM.
	ctx.SetVMConfig(2, 512)

	// Use a host directory as the root filesystem.
	ctx.SetRoot("/path/to/rootfs")

	// Set the command to run inside the VM.
	ctx.SetExec("/bin/uname", []string{"/bin/uname", "-a"}, nil)

	// Start the VM. Does not return on success.
	if err := ctx.StartEnter(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
```

## API Overview

The typical workflow is:

1. Create a context with `krun.CreateContext()`
2. Configure the VM using methods on `*krun.Context`
3. Call `ctx.StartEnter()` to launch the microVM

### Package-level functions

| Function | Description |
|----------|-------------|
| `CreateContext()` | Create a new VM configuration context |
| `SetLogLevel(level)` | Set library log verbosity |
| `InitLog(...)` | Initialize logging with full control |
| `HasFeature(feature)` | Check if a feature was enabled at build time |
| `GetMaxVCPUs()` | Query max vCPUs supported by the hypervisor |
| `CheckNestedVirt()` | Check nested virtualization support (macOS) |

### Context methods

| Category | Methods |
|----------|---------|
| VM config | `SetVMConfig`, `SetRoot`, `SetNestedVirt`, `SetMappedVolumes` |
| Execution | `SetExec`, `SetWorkdir`, `SetEnv`, `SetRlimits` |
| Disks | `AddDisk`, `AddDisk2`, `AddDisk3`, `SetRootDiskRemount` |
| Filesystem | `AddVirtioFS` |
| Network | `SetNetUnix`, `SetNetGram`, `SetNetTap`, `SetNetMAC`, `AddPortMap` |
| GPU | `SetGPUOptions`, `SetGPUOptions2`, `SetDisplayOptions`, `SetSoundOptions` |
| Console | `AddVirtioConsole`, `AddVirtioConsoleDefault`, `AddSerialConsole` |
| Vsock | `AddVsockPort`, `SetTSI`, `AddTSIPortMap` |
| Kernel | `SetFirmware`, `SetKernel` |
| Lifecycle | `StartEnter`, `Free` |

### Error handling

Errors are returned as `*krun.Error` which wraps a `syscall.Errno`, so you can use `errors.Is`:

```go
if errors.Is(err, syscall.EINVAL) {
	// handle invalid argument
}
```

## Examples

See the [`examples/`](examples/) directory:

- **[features](examples/features/)** — Query library capabilities (no rootfs needed)
- **[basic](examples/basic/)** — Run a command in a microVM using a host directory
- **[vm-with-disk](examples/vm-with-disk/)** — Boot from a disk image with a custom kernel

## License

Apache 2.0 — see [LICENSE](LICENSE) for details.
