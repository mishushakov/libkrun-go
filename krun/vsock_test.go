package krun

import (
	"path/filepath"
	"testing"
)

func TestAddVsockPort(t *testing.T) {
	ctx := newTestContext(t)
	sockPath := filepath.Join(t.TempDir(), "vsock.sock")
	if err := ctx.AddVsockPort(VsockPortConfig{Port: 5000, Path: sockPath}); err != nil {
		t.Fatal(err)
	}
}

func TestAddVsockPort_Listen(t *testing.T) {
	ctx := newTestContext(t)

	t.Run("listen_false", func(t *testing.T) {
		sockPath := filepath.Join(t.TempDir(), "vsock.sock")
		if err := ctx.AddVsockPort(VsockPortConfig{Port: 5001, Path: sockPath}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("listen_true", func(t *testing.T) {
		sockPath := filepath.Join(t.TempDir(), "vsock2.sock")
		if err := ctx.AddVsockPort(VsockPortConfig{Port: 5002, Path: sockPath, Listen: true}); err != nil {
			t.Fatal(err)
		}
	})
}

func TestAddVsock(t *testing.T) {
	ctx := newTestContext(t)
	ctx.DisableImplicitVsock()
	if err := ctx.AddVsock(0); err != nil {
		t.Fatal(err)
	}
}

func TestDisableImplicitVsock(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.DisableImplicitVsock(); err != nil {
		t.Fatal(err)
	}
}
