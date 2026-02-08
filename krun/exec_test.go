package krun

import "testing"

func TestSetWorkdir(t *testing.T) {
	ctx := newTestContext(t)
	if err := ctx.SetWorkdir("/tmp"); err != nil {
		t.Fatal(err)
	}
}

func TestSetExec(t *testing.T) {
	ctx := newTestContext(t)

	t.Run("with_explicit_env", func(t *testing.T) {
		err := ctx.SetExec("/bin/sh", []string{"sh", "-c", "echo hello"}, []string{"PATH=/usr/bin"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with_nil_env", func(t *testing.T) {
		err := ctx.SetExec("/bin/sh", []string{"sh"}, nil)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with_empty_env", func(t *testing.T) {
		err := ctx.SetExec("/bin/sh", []string{"sh"}, []string{})
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestSetEnv(t *testing.T) {
	ctx := newTestContext(t)

	t.Run("explicit", func(t *testing.T) {
		err := ctx.SetEnv([]string{"FOO=bar", "BAZ=qux"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("nil_auto", func(t *testing.T) {
		err := ctx.SetEnv(nil)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestSetRlimits(t *testing.T) {
	ctx := newTestContext(t)
	err := ctx.SetRlimits([]string{"RLIMIT_NOFILE=1024:4096"})
	if err != nil {
		t.Fatal(err)
	}
}
