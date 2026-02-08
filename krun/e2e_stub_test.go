//go:build !linux

package krun

// e2eHelper is a no-op stub on non-Linux platforms.
// The real implementation is in e2e_test.go (linux only).
func e2eHelper() {}
