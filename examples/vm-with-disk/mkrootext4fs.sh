#!/usr/bin/env bash
set -euo pipefail

usage() {
	echo "Usage: $0 <docker-image> [output-file] [size-mb]" >&2
	echo "" >&2
	echo "Create an ext4 root filesystem image from a Docker image." >&2
	echo "" >&2
	echo "Arguments:" >&2
	echo "  docker-image  Docker image to use (e.g. alpine, ubuntu:22.04)" >&2
	echo "  output-file   Output image path (default: ./rootfs.ext4)" >&2
	echo "  size-mb       Image size in MiB (default: 1024)" >&2
	echo "" >&2
	echo "Examples:" >&2
	echo "  $0 alpine" >&2
	echo "  $0 ubuntu:22.04 ./rootfs.ext4 2048" >&2
	exit 1
}

[ $# -lt 1 ] && usage

image="$1"
outfile="${2:-./rootfs.ext4}"
size_mb="${3:-1024}"

if [ -f "$outfile" ]; then
	echo "error: output file '$outfile' already exists" >&2
	exit 1
fi

# Check for required tools.
for cmd in docker dd mkfs.ext4; do
	if ! command -v "$cmd" >/dev/null 2>&1; then
		echo "error: '$cmd' is required but not found" >&2
		exit 1
	fi
done

# macOS uses a different mount approach than Linux.
if [ "$(uname)" = "Darwin" ]; then
	echo "error: this script requires Linux (for loop mounting)" >&2
	echo "hint: run this inside a Linux VM or container" >&2
	exit 1
fi

echo "Pulling $image..."
docker pull "$image"

echo "Creating container..."
cid=$(docker create "$image")
trap 'docker rm "$cid" >/dev/null' EXIT

echo "Creating ${size_mb}M ext4 image at $outfile..."
dd if=/dev/zero of="$outfile" bs=1M count="$size_mb" status=progress
mkfs.ext4 -q "$outfile"

mountdir=$(mktemp -d)
trap 'sudo umount "$mountdir" 2>/dev/null; rmdir "$mountdir"; docker rm "$cid" >/dev/null' EXIT

echo "Mounting image..."
sudo mount -o loop "$outfile" "$mountdir"

echo "Exporting filesystem..."
docker export "$cid" | sudo tar -x -C "$mountdir"

echo "Unmounting..."
sudo umount "$mountdir"
rmdir "$mountdir"

echo "Done: $outfile"
