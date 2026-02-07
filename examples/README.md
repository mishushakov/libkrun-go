# Examples

## Prerequisites

You need libkrun and libkrunfw installed on your system.

### libkrunfw (firmware)

libkrunfw bundles a Linux kernel used by the microVM. On macOS it is available via Homebrew:

```bash
brew install slp/krun/libkrunfw
```

If libkrun was installed to `/usr/local/lib` but Homebrew puts libkrunfw in `/opt/homebrew`, symlink it:

```bash
sudo ln -sf /opt/homebrew/opt/libkrunfw/lib/libkrunfw.5.dylib /usr/local/lib/
```

### Building libkrun from source

```bash
cd libkrun
make
sudo make install
```

This installs `libkrun.dylib` (macOS) or `libkrun.so` (Linux) into `/usr/local/lib` and the header into `/usr/local/include`.

If you install to a non-standard location, set the linker search path:

```bash
# macOS
export DYLD_LIBRARY_PATH=/path/to/lib:$DYLD_LIBRARY_PATH

# Linux
export LD_LIBRARY_PATH=/path/to/lib:$LD_LIBRARY_PATH
```

### macOS: Hypervisor entitlement

On macOS, any binary that uses Hypervisor.framework must be signed with the `com.apple.security.hypervisor` entitlement. This means **`go run` will not work** for examples that start a VM — you must build, sign, then run:

```bash
go build -o basic .
codesign --entitlements entitlements.plist --force -s - basic
./basic /path/to/rootfs /bin/uname -a
```

An `entitlements.plist` is provided in the `basic/` example directory. You can reuse it for other examples.

## Build tags

Some libkrun features are optional and may not be present in every build. The Go bindings use build tags to match your libkrun's feature set:

| Build tag | libkrun feature | Functions |
|-----------|----------------|-----------|
| `krun_blk` | `BLK` | `AddDisk`, `AddDisk2`, `AddDisk3`, `SetRootDiskRemount` |
| `krun_net` | `NET` | `AddNetUnixStream`, `AddNetUnixGram`, `AddNetTap`, `SetNetMac` |
| `krun_tee` | `TEE` | `SetTEEConfigFile` |

Without the corresponding tag, these functions are still available but return `ENOSYS` at runtime.

To enable optional features, pass `-tags` when building:

```bash
# Enable block device support
go build -tags krun_blk ./...

# Enable multiple features
go build -tags "krun_blk,krun_net" ./...
```

You can check which features your libkrun was built with by running the `features` example.

## Running the examples

### features — Query library capabilities

The easiest example to start with — doesn't require a rootfs or disk image.

```bash
cd examples/features
go run .
```

Sample output:

```
Max vCPUs: 8
Nested virtualization: true

Compile-time features:
  Networking                yes
  Block devices             yes
  GPU                       no
  Sound                     no
  Input                     no
  EFI                       yes
  TEE                       no
  AMD SEV                   no
  Intel TDX                 no
  AWS Nitro                 no
  Virgl Resource Map2       no
```

### basic — Run a command in a microVM

Runs a command inside a microVM using a host directory as the root filesystem.

```bash
cd examples/basic
go build -o basic .
codesign --entitlements entitlements.plist --force -s - basic
./basic /path/to/rootfs /bin/uname -a
```

The rootfs directory should contain a minimal Linux filesystem (with `/bin`, `/lib`, etc.). The easiest way is to use the included helper script to extract one from a Docker image:

```bash
# Create a rootfs from an Alpine image
../mkrootfs.sh alpine ./rootfs

# Or from Ubuntu
../mkrootfs.sh ubuntu:22.04 ./rootfs
```

You can also create one manually with `debootstrap`:

```bash
sudo debootstrap --variant=minbase bookworm ./rootfs
```

### vm-with-disk — Boot from a disk image

Boots a microVM from a disk image with a custom kernel and optional virtio-fs shared directory.

This example requires the `krun_blk` build tag since it uses disk images. On macOS, you must build and sign the binary (see [macOS: Hypervisor entitlement](#macos-hypervisor-entitlement)):

```bash
cd examples/vm-with-disk
go build -tags krun_blk -o vm-with-disk .
codesign --entitlements ../basic/entitlements.plist --force -s - vm-with-disk
./vm-with-disk \
  -kernel /path/to/vmlinux \
  -disk /path/to/rootfs.ext4
```

All flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-kernel` | *(required)* | Path to kernel image |
| `-disk` | *(required)* | Path to root disk image |
| `-format` | `raw` | Disk format: `raw`, `qcow2`, `vmdk` |
| `-shared` | | Host directory to share via virtio-fs |
| `-vcpus` | `2` | Number of vCPUs |
| `-ram` | `1024` | RAM in MiB |
| `-cmdline` | `console=hvc0 root=/dev/vda1 rw` | Kernel command line |

Example with a shared directory:

```bash
./vm-with-disk \
  -kernel /path/to/vmlinux \
  -disk /path/to/rootfs.ext4 \
  -shared /home/user/workspace \
  -vcpus 4 \
  -ram 2048
```

Inside the guest, mount the shared directory:

```bash
mount -t virtiofs shared /mnt
```

#### Creating a disk image

```bash
# Create a 1GB raw disk image with ext4
dd if=/dev/zero of=rootfs.img bs=1M count=1024
mkfs.ext4 rootfs.img

# Mount and populate it
sudo mount -o loop rootfs.img /mnt
sudo debootstrap --variant=minbase bookworm /mnt
sudo umount /mnt
```

## Troubleshooting

**`Couldn't find or load libkrunfw.5.dylib`** — libkrunfw is not installed or not in the library search path. Install it via Homebrew and symlink if needed (see [Prerequisites](#prerequisites)).

**`VmSetup(VmCreate)`** (macOS) — The binary is missing the Hypervisor entitlement. You cannot use `go run` for examples that start a VM. Build and sign the binary instead (see [macOS: Hypervisor entitlement](#macos-hypervisor-entitlement)).

**`TooLarge` panic on aarch64** — The kernel command line exceeded the 2048-byte limit on aarch64. This typically happens when passing `nil` for the environment in `SetExec`, which inherits the full host environment. Pass an explicit minimal environment instead.

**`Undefined symbols for architecture`** — Your libkrun was built without some optional features. Use build tags to match your build (see [Build tags](#build-tags) above), or rebuild libkrun with the needed features:

```bash
cd libkrun
make BLK=1 NET=1
sudo make install
```

**`library not found for -lkrun`** — libkrun is not installed or not in the linker search path. Build and install it, or set `CGO_LDFLAGS`:

```bash
CGO_LDFLAGS="-L/path/to/libkrun/target/release" go build .
```

**`dyld: Library not loaded: libkrun.dylib`** (macOS runtime) — The dynamic library isn't in the runtime search path:

```bash
export DYLD_LIBRARY_PATH=/usr/local/lib
```

**`permission denied`** (Linux) — KVM access requires the user to be in the `kvm` group:

```bash
sudo usermod -aG kvm $USER
# Log out and back in
```
