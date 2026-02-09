# libkrun-go

Go bindings for [libkrun](https://github.com/containers/libkrun), a dynamic library for creating lightweight microVMs using KVM (Linux) or HVF (macOS/ARM64).

## Installation

```bash
go get github.com/mishushakov/libkrun-go/krun
```

libkrun must be installed on your system. The bundled header in `libkrun/include/libkrun.h` is used at build time when the submodule is present. Otherwise, `pkg-config` is used to locate headers and the shared library (`libkrun.so` or `libkrun.dylib`), and you can override paths with `CGO_CFLAGS`/`CGO_LDFLAGS`.

### Building libkrun from source

```bash
cd libkrun
make
sudo make install
```

### macOS

libkrunfw (the firmware) is available via Homebrew (repo [here](https://github.com/slp/homebrew-krun)):

```bash
brew tap slp/krun
brew install libkrun pkg-config
```

libkrun itself must be installed (either via Homebrew or from source).

`pkg-config` is used on macOS to locate libkrun headers and libraries from Homebrew.
If you installed libkrun via Homebrew, make sure `pkg-config` can find `libkrun.pc` (Homebrew usually handles this automatically). If not, set:

```bash
export PKG_CONFIG_PATH="/opt/homebrew/lib/pkgconfig:$PKG_CONFIG_PATH"
```

On macOS, binaries are linked with an rpath to `/opt/homebrew/lib` for a smoother Homebrew experience. If the runtime loader still cannot find `libkrunfw.5.dylib`, set a library search path before running your binary:

```bash
export DYLD_LIBRARY_PATH="/opt/homebrew/lib:$DYLD_LIBRARY_PATH"
```

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/mishushakov/libkrun-go/krun"
)

func main() {
	// Create a new VM configuration context.
	ctx, err := krun.CreateContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Configure 2 vCPUs and 512 MiB of RAM.
	ctx.SetVMConfig(krun.VMConfig{NumVCPUs: 2, RAMMiB: 512})

	// Use a host directory as the root filesystem.
	ctx.SetRoot("/path/to/rootfs")

	// Set the command to run inside the VM.
	ctx.SetExec(krun.ExecConfig{
		Path: "/bin/uname",
		Args: []string{"/bin/uname", "-a"},
	})

	// Start the VM. Does not return on success.
	if err := ctx.StartEnter(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
```

## Build Tags

Some libkrun features are optional and gated behind Go build tags. Without the corresponding tag, calls to those functions return `syscall.ENOSYS`.

| Tag | Feature |
|-----|---------|
| `krun_blk` | Block device / disk support (`AddDisk`, `SetRootDiskRemount`) |
| `krun_net` | Network backends (`AddNetUnixStream`, `AddNetUnixGram`, `AddNetTap`, `SetNetMac`) |
| `krun_tee` | TEE configuration (`SetTEEConfigFile`) |

Build with tags:

```bash
go build -tags "krun_blk,krun_net" ./...
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
| `InitLog(targetFD, level, style, options)` | Initialize logging with full control |
| `HasFeature(feature)` | Check if a feature was enabled at build time |
| `GetMaxVCPUs()` | Query max vCPUs supported by the hypervisor |
| `CheckNestedVirt()` | Check nested virtualization support (macOS) |

### Context methods

#### VM configuration

| Method | Description |
|--------|-------------|
| `SetVMConfig(VMConfig)` | Set vCPU count and RAM |
| `SetRoot(rootPath)` | Set root filesystem path |
| `SetNestedVirt(enabled)` | Enable/disable nested virtualization (macOS) |
| `SplitIRQChip(enable)` | Split IRQCHIP between host and guest |
| `SetUID(uid)` | Set user ID before VM startup |
| `SetGID(gid)` | Set group ID before VM startup |
| `SetSMBIOSOEMStrings(oemStrings)` | Set SMBIOS OEM Strings |
| `GetShutdownEventFD()` | Get file descriptor for shutdown signaling (libkrun-efi) |

#### Execution

| Method | Description |
|--------|-------------|
| `SetExec(ExecConfig)` | Set executable, args, and environment |
| `SetWorkdir(workdirPath)` | Set working directory for the executable |
| `SetEnv(envp)` | Set environment variables |
| `SetRlimits(rlimits)` | Set guest resource limits |

#### Disks (requires `krun_blk` tag)

| Method | Description |
|--------|-------------|
| `AddDisk(DiskConfig)` | Add a disk image with full options |
| `SetRootDiskRemount(RootDiskRemountConfig)` | Mount a block device as root filesystem |

#### Filesystem

| Method | Description |
|--------|-------------|
| `AddVirtioFS(VirtioFSConfig)` | Add a virtio-fs shared directory |

#### Network

| Method | Description | Tag |
|--------|-------------|-----|
| `SetPortMap(portMap)` | Configure host-to-guest TCP port mappings | — |
| `AddNetUnixStream(NetUnixConfig)` | Add net device via unix stream (e.g., passt) | `krun_net` |
| `AddNetUnixGram(NetUnixConfig)` | Add net device via unix dgram (e.g., gvproxy) | `krun_net` |
| `AddNetTap(NetTapConfig)` | Add net device via TAP | `krun_net` |
| `SetNetMac(mac)` | Set MAC address for passt backend | `krun_net` |

#### GPU, display, input, and sound

| Method | Description |
|--------|-------------|
| `SetGPUOptions(GPUConfig)` | Enable and configure virtio-gpu |
| `AddDisplay(DisplayConfig)` | Add a display output (returns display ID) |
| `DisplaySetEDID(displayID, edidBlob)` | Set custom EDID for a display |
| `DisplaySetDPI(displayID, dpi)` | Set display DPI |
| `DisplaySetPhysicalSize(displayID, widthMM, heightMM)` | Set physical display dimensions |
| `DisplaySetRefreshRate(displayID, refreshRate)` | Set display refresh rate |
| `SetDisplayBackend(backend, size)` | Set display backend |
| `AddInputDeviceFD(inputFD)` | Passthrough a host input device |
| `AddInputDevice(configBackend, configSize, eventsBackend, eventsSize)` | Add input device with custom backends |
| `SetSndDevice(enable)` | Enable/disable virtio-snd |

#### Console and serial

| Method | Description |
|--------|-------------|
| `SetConsoleOutput(filepath)` | Redirect implicit console output to a file |
| `DisableImplicitConsole()` | Disable the implicit console device |
| `SetKernelConsole(consoleID)` | Set kernel `console=` parameter |
| `AddVirtioConsoleDefault(VirtioConsoleConfig)` | Add virtio-console with automatic detection |
| `AddVirtioConsoleMultiport()` | Create multi-port virtio-console (returns ID) |
| `AddConsolePortTTY(ConsolePortTTYConfig)` | Add TTY port to multi-port console |
| `AddConsolePortInOut(ConsolePortInOutConfig)` | Add generic I/O port to multi-port console |
| `AddSerialConsoleDefault(SerialConsoleConfig)` | Add legacy serial device |

#### Vsock

| Method | Description |
|--------|-------------|
| `AddVsockPort(VsockPortConfig)` | Map vsock port to a host UNIX socket |
| `AddVsock(tsiFeatures)` | Add vsock device with TSI features |
| `DisableImplicitVsock()` | Disable the default vsock device |

#### Kernel and firmware

| Method | Description |
|--------|-------------|
| `SetFirmware(firmwarePath)` | Load firmware into the microVM |
| `SetKernel(KernelConfig)` | Load kernel with initramfs and command line |

#### TEE (requires `krun_tee` tag)

| Method | Description |
|--------|-------------|
| `SetTEEConfigFile(filepath)` | Set TEE configuration file (libkrun-sev) |

#### Lifecycle

| Method | Description |
|--------|-------------|
| `ID()` | Get the underlying context ID |
| `StartEnter()` | Start and enter the microVM (does not return on success) |
| `Free()` | Release the configuration context |

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
