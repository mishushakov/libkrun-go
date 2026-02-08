package krun

import (
	"path/filepath"
	"testing"
)

func TestAddVsockPort(t *testing.T) {
	ctx := newTestContext(t)
	sockPath := filepath.Join(t.TempDir(), "vsock.sock")
	if err := ctx.AddVsockPort(5000, sockPath); err != nil {
		t.Fatal(err)
	}
}

func TestAddVsockPort2(t *testing.T) {
	ctx := newTestContext(t)
	sockPath := filepath.Join(t.TempDir(), "vsock.sock")

	t.Run("listen_false", func(t *testing.T) {
		if err := ctx.AddVsockPort2(5001, sockPath, false); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("listen_true", func(t *testing.T) {
		sockPath2 := filepath.Join(t.TempDir(), "vsock2.sock")
		if err := ctx.AddVsockPort2(5002, sockPath2, true); err != nil {
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
