package krun

import (
	"errors"
	"syscall"
	"testing"
)

func TestSetPortMap(t *testing.T) {
	ctx := newTestContext(t)

	t.Run("explicit_ports", func(t *testing.T) {
		err := ctx.SetPortMap([]string{"8080:80", "4433:443"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty_expose_none", func(t *testing.T) {
		err := ctx.SetPortMap([]string{})
		if err != nil {
			t.Fatal(err)
		}
	})
}

// Stub tests — these verify behavior when built without -tags krun_net.
// When built with the tag, these tests are still valid (the real implementations
// may succeed or fail differently).

func TestAddNetUnixStream_Stub(t *testing.T) {
	ctx := newTestContext(t)
	mac := [6]byte{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x01}
	err := ctx.AddNetUnixStream("/tmp/test.sock", -1, mac, 0, 0)
	if err == nil {
		return // built with krun_net tag, real impl succeeded
	}
	if !errors.Is(err, syscall.ENOSYS) {
		t.Fatalf("expected ENOSYS, got %v", err)
	}
}

func TestAddNetUnixGram_Stub(t *testing.T) {
	ctx := newTestContext(t)
	mac := [6]byte{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x02}
	err := ctx.AddNetUnixGram("/tmp/test.sock", -1, mac, 0, 0)
	if err == nil {
		return
	}
	if !errors.Is(err, syscall.ENOSYS) {
		t.Fatalf("expected ENOSYS, got %v", err)
	}
}

func TestAddNetTap_Stub(t *testing.T) {
	ctx := newTestContext(t)
	mac := [6]byte{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x03}
	err := ctx.AddNetTap("tap0", mac, 0, 0)
	if err == nil {
		return
	}
	if !errors.Is(err, syscall.ENOSYS) {
		t.Fatalf("expected ENOSYS, got %v", err)
	}
}

func TestSetNetMac_Stub(t *testing.T) {
	ctx := newTestContext(t)
	mac := [6]byte{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x04}
	err := ctx.SetNetMac(mac)
	if err == nil {
		return
	}
	if !errors.Is(err, syscall.ENOSYS) {
		t.Fatalf("expected ENOSYS, got %v", err)
	}
}
