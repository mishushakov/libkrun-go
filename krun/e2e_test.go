//go:build linux

package krun

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// e2eHelper is called by TestMain when KRUN_E2E_HELPER=1.
// It configures a minimal VM and calls StartEnter, which never returns
// on success (the VMM calls exit with the guest's exit code).
func e2eHelper() {
	rootfs := os.Getenv("KRUN_E2E_ROOTFS")
	execPath := os.Getenv("KRUN_E2E_EXEC")
	if rootfs == "" || execPath == "" {
		fmt.Fprintln(os.Stderr, "e2eHelper: KRUN_E2E_ROOTFS and KRUN_E2E_EXEC must be set")
		os.Exit(1)
	}

	if err := SetLogLevel(LogLevelOff); err != nil {
		fmt.Fprintf(os.Stderr, "e2eHelper: SetLogLevel: %v\n", err)
		os.Exit(1)
	}

	ctx, err := CreateContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "e2eHelper: CreateContext: %v\n", err)
		os.Exit(1)
	}

	if err := ctx.SetVMConfig(VMConfig{NumVCPUs: 1, RAMMiB: 256}); err != nil {
		fmt.Fprintf(os.Stderr, "e2eHelper: SetVMConfig: %v\n", err)
		os.Exit(1)
	}

	if err := ctx.SetRoot(rootfs); err != nil {
		fmt.Fprintf(os.Stderr, "e2eHelper: SetRoot: %v\n", err)
		os.Exit(1)
	}

	if err := ctx.SetExec(ExecConfig{Path: execPath, Args: []string{execPath}, Env: []string{}}); err != nil {
		fmt.Fprintf(os.Stderr, "e2eHelper: SetExec: %v\n", err)
		os.Exit(1)
	}

	if err := ctx.SetWorkdir("/"); err != nil {
		fmt.Fprintf(os.Stderr, "e2eHelper: SetWorkdir: %v\n", err)
		os.Exit(1)
	}

	// StartEnter consumes the context and never returns on success.
	if err := ctx.StartEnter(); err != nil {
		fmt.Fprintf(os.Stderr, "e2eHelper: StartEnter: %v\n", err)
		os.Exit(1)
	}
}

// skipIfNoKVM skips the test if /dev/kvm is not accessible.
func skipIfNoKVM(t *testing.T) {
	t.Helper()
	if _, err := os.Stat("/dev/kvm"); err != nil {
		t.Skip("skipping: /dev/kvm not available")
	}
}

// buildStaticGuest compiles a tiny static C program into dir/name.
// Returns the absolute path to the compiled binary.
func buildStaticGuest(t *testing.T, dir, name, cSource string) string {
	t.Helper()
	srcPath := filepath.Join(dir, name+".c")
	binPath := filepath.Join(dir, name)

	if err := os.WriteFile(srcPath, []byte(cSource), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command("cc", "-static", "-o", binPath, srcPath).CombinedOutput()
	if err != nil {
		t.Fatalf("cc -static failed: %v\n%s", err, out)
	}
	return binPath
}

// runE2E re-executes the test binary as a child process with KRUN_E2E_HELPER=1.
// Returns the child's combined stdout/stderr and exit code.
func runE2E(t *testing.T, rootfs, execPath string) (stdout string, exitCode int) {
	t.Helper()
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(self)
	cmd.Env = append(os.Environ(),
		"KRUN_E2E_HELPER=1",
		"KRUN_E2E_ROOTFS="+rootfs,
		"KRUN_E2E_EXEC="+execPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(out), exitErr.ExitCode()
		}
		t.Fatalf("failed to run e2e helper: %v\n%s", err, out)
	}
	return string(out), 0
}

// TestE2EBootAndExit boots a VM with a guest that writes "OK" to stdout
// and exits 0, then verifies both.
func TestE2EBootAndExit(t *testing.T) {
	skipIfNoKVM(t)

	rootfs := t.TempDir()

	binPath := buildStaticGuest(t, rootfs, "guest", `
#include <unistd.h>
int main(void) {
    write(1, "OK\n", 3);
    return 0;
}
`)

	// execPath is relative to rootfs
	relPath := "/" + filepath.Base(binPath)

	stdout, exitCode := runE2E(t, rootfs, relPath)
	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0\noutput: %s", exitCode, stdout)
	}
	if !strings.Contains(stdout, "OK") {
		t.Errorf("stdout = %q, want it to contain %q", stdout, "OK")
	}
}

// TestE2ENonZeroExit boots a VM with a guest that exits 42 and verifies
// the exit code is propagated.
func TestE2ENonZeroExit(t *testing.T) {
	skipIfNoKVM(t)

	rootfs := t.TempDir()

	buildStaticGuest(t, rootfs, "guest", `
int main(void) {
    return 42;
}
`)

	stdout, exitCode := runE2E(t, rootfs, "/guest")
	if exitCode != 42 {
		t.Errorf("exit code = %d, want 42\noutput: %s", exitCode, stdout)
	}
}
