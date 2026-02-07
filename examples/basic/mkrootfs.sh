#!/usr/bin/env bash
set -euo pipefail

usage() {
	echo "Usage: $0 <docker-image> [output-dir]" >&2
	echo "" >&2
	echo "Create a root filesystem directory from a Docker image." >&2
	echo "" >&2
	echo "Examples:" >&2
	echo "  $0 alpine" >&2
	echo "  $0 ubuntu:22.04 ./rootfs" >&2
	exit 1
}

[ $# -lt 1 ] && usage

image="$1"
outdir="${2:-./rootfs}"

if [ -d "$outdir" ]; then
	echo "error: output directory '$outdir' already exists" >&2
	exit 1
fi

echo "Pulling $image..."
docker pull "$image"

echo "Creating container..."
cid=$(docker create "$image")
trap 'docker rm "$cid" >/dev/null' EXIT

echo "Exporting filesystem to $outdir..."
mkdir -p "$outdir"
docker export "$cid" | tar -x -C "$outdir"

echo "Done: $outdir"
