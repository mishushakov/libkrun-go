package krun

import (
	"errors"
	"syscall"
	"testing"
)

func TestSetGPUOptions(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.SetGPUOptions(GPUConfig{VirglFlags: VirglUseSurfaceless | VirglUseEGL}); err != nil {
		t.Fatal(err)
	}
}

func TestSetGPUOptions_WithShmSize(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.SetGPUOptions(GPUConfig{VirglFlags: VirglUseSurfaceless, ShmSize: 256 * 1024 * 1024}); err != nil {
		t.Fatal(err)
	}
}

func TestSetGPURenderServerFD(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.SetGPURenderServerFD(3); err != nil {
		t.Fatal(err)
	}
	if err := ctx.SetGPURenderServerFD(-1); !errors.Is(err, syscall.EINVAL) {
		t.Errorf("SetGPURenderServerFD(-1) = %v, want EINVAL", err)
	}
}

func TestSetSndDevice(t *testing.T) {
	ctx := newTestContext(t)
	for _, enable := range []bool{true, false} {
		if err := ctx.SetSndDevice(enable); err != nil {
			t.Errorf("SetSndDevice(%v) = %v", enable, err)
		}
	}
}
