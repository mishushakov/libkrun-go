package krun

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"unsafe"
)

func TestMain(m *testing.M) {
	// SetLogLevel must be called before any other libkrun function
	// because the Rust env_logger can only be initialized once.
	if err := SetLogLevel(LogLevelOff); err != nil {
		fmt.Fprintf(os.Stderr, "SetLogLevel: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func newTestContext(t *testing.T) *Context {
	t.Helper()
	ctx, err := CreateContext()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ctx.Free() })
	return ctx
}

// Error type tests

func TestError_Format(t *testing.T) {
	err := &Error{Func: "krun_test_func", Errno: syscall.EINVAL}
	got := err.Error()
	want := "krun: krun_test_func: invalid argument"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestError_Unwrap(t *testing.T) {
	err := &Error{Func: "krun_test_func", Errno: syscall.EPERM}
	unwrapped := err.Unwrap()
	if unwrapped != syscall.EPERM {
		t.Errorf("Unwrap() = %v, want EPERM", unwrapped)
	}
}

func TestError_Is(t *testing.T) {
	err := &Error{Func: "krun_test_func", Errno: syscall.EINVAL}
	if !errors.Is(err, syscall.EINVAL) {
		t.Error("errors.Is(err, EINVAL) = false, want true")
	}
	if errors.Is(err, syscall.EPERM) {
		t.Error("errors.Is(err, EPERM) = true, want false")
	}
}

func TestError_As(t *testing.T) {
	var orig error = &Error{Func: "krun_test_func", Errno: syscall.ENOENT}
	var kErr *Error
	if !errors.As(orig, &kErr) {
		t.Fatal("errors.As failed")
	}
	if kErr.Func != "krun_test_func" {
		t.Errorf("Func = %q, want %q", kErr.Func, "krun_test_func")
	}
	if kErr.Errno != syscall.ENOENT {
		t.Errorf("Errno = %v, want ENOENT", kErr.Errno)
	}
}

// checkRet / retError tests

func TestCheckRet_Success(t *testing.T) {
	if err := checkRet(0, "test"); err != nil {
		t.Errorf("checkRet(0) = %v, want nil", err)
	}
}

func TestCheckRet_PositiveSuccess(t *testing.T) {
	if err := checkRet(42, "test"); err != nil {
		t.Errorf("checkRet(42) = %v, want nil", err)
	}
}

func TestCheckRet_Error(t *testing.T) {
	err := checkRet(-22, "krun_test") // -EINVAL
	if err == nil {
		t.Fatal("checkRet(-22) = nil, want error")
	}
	if !errors.Is(err, syscall.EINVAL) {
		t.Errorf("expected EINVAL, got %v", err)
	}
}

func TestRetError(t *testing.T) {
	err := retError(-1, "krun_test") // -EPERM
	if err == nil {
		t.Fatal("retError(-1) = nil, want error")
	}
	if !errors.Is(err, syscall.EPERM) {
		t.Errorf("expected EPERM, got %v", err)
	}
}

// stringsToCArray / freeCStringArray tests

func TestStringsToCArray_Nil(t *testing.T) {
	arr := stringsToCArray(nil)
	if arr != nil {
		t.Error("stringsToCArray(nil) != nil")
	}
	// freeCStringArray should handle nil safely
	freeCStringArray(nil, 0)
}

func TestStringsToCArray_Empty(t *testing.T) {
	arr := stringsToCArray([]string{})
	if arr == nil {
		t.Fatal("stringsToCArray([]) = nil, want non-nil")
	}
	defer freeCStringArray(arr, 0)
	// First element should be nil (null terminator)
	slice := unsafe.Slice(arr, 1)
	if slice[0] != nil {
		t.Error("empty array first element should be nil")
	}
}

func TestStringsToCArray_Values(t *testing.T) {
	strs := []string{"hello", "world"}
	arr := stringsToCArray(strs)
	if arr == nil {
		t.Fatal("stringsToCArray returned nil")
	}
	defer freeCStringArray(arr, len(strs))
	// Verify null termination
	slice := unsafe.Slice(arr, len(strs)+1)
	if slice[len(strs)] != nil {
		t.Error("array should be null-terminated")
	}
}

// Context lifecycle tests

func TestCreateContext(t *testing.T) {
	ctx, err := CreateContext()
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Free()
	// ID should be a small non-negative value
	t.Logf("context ID: %d", ctx.ID())
}

func TestContextFree(t *testing.T) {
	ctx, err := CreateContext()
	if err != nil {
		t.Fatal(err)
	}
	if err := ctx.Free(); err != nil {
		t.Fatal(err)
	}
}

// Package-level function tests

func TestHasFeature(t *testing.T) {
	features := []struct {
		name    string
		feature Feature
	}{
		{"Net", FeatureNet},
		{"BLK", FeatureBLK},
		{"GPU", FeatureGPU},
		{"SND", FeatureSND},
		{"Input", FeatureInput},
		{"EFI", FeatureEFI},
		{"TEE", FeatureTEE},
	}
	for _, f := range features {
		t.Run(f.name, func(t *testing.T) {
			supported, err := HasFeature(f.feature)
			if err != nil {
				t.Errorf("HasFeature(%s) error: %v", f.name, err)
			}
			t.Logf("Feature %s supported: %v", f.name, supported)
		})
	}
}

func TestGetMaxVCPUs(t *testing.T) {
	max, err := GetMaxVCPUs()
	if err != nil {
		t.Fatal(err)
	}
	if max <= 0 {
		t.Errorf("GetMaxVCPUs() = %d, want > 0", max)
	}
	t.Logf("max vCPUs: %d", max)
}

func TestCheckNestedVirt(t *testing.T) {
	supported, err := CheckNestedVirt()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("nested virt supported: %v", supported)
}
