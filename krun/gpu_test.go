package krun

import "testing"

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

func TestSetSndDevice(t *testing.T) {
	ctx := newTestContext(t)
	for _, enable := range []bool{true, false} {
		if err := ctx.SetSndDevice(enable); err != nil {
			t.Errorf("SetSndDevice(%v) = %v", enable, err)
		}
	}
}
