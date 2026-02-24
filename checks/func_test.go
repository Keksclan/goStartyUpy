package checks

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestNew_NilFnFails(t *testing.T) {
	c := New("nilcheck", nil)
	res := c.Run(t.Context())
	if res.OK {
		t.Fatal("expected FAIL for nil fn")
	}
	if res.Error != "nil check function" {
		t.Errorf("unexpected error: %q", res.Error)
	}
}

func TestNew_PanicConvertedToFail(t *testing.T) {
	c := New("panicker", func(_ context.Context) error {
		panic("boom")
	})
	res := c.Run(t.Context())
	if res.OK {
		t.Fatal("expected FAIL after panic")
	}
	if !strings.HasPrefix(res.Error, "panic: ") {
		t.Errorf("expected error starting with 'panic: ', got %q", res.Error)
	}
	if !strings.Contains(res.Error, "boom") {
		t.Errorf("expected panic message to contain 'boom', got %q", res.Error)
	}
}

func TestNew_EmptyLabelDefaultsToCustom(t *testing.T) {
	c := New("", func(_ context.Context) error { return nil })
	if c.Name() != "custom" {
		t.Errorf("expected name 'custom', got %q", c.Name())
	}
}

func TestNew_SuccessfulRun(t *testing.T) {
	c := New("ok-check", func(_ context.Context) error { return nil })
	res := c.Run(t.Context())
	if !res.OK {
		t.Fatal("expected OK")
	}
	if res.Error != "" {
		t.Errorf("unexpected error: %q", res.Error)
	}
	if res.Duration < 0 {
		t.Error("expected non-negative duration")
	}
}

func TestNew_ErrorRun(t *testing.T) {
	c := New("err-check", func(_ context.Context) error {
		return errors.New("something broke")
	})
	res := c.Run(t.Context())
	if res.OK {
		t.Fatal("expected FAIL")
	}
	if res.Error != "something broke" {
		t.Errorf("unexpected error: %q", res.Error)
	}
}

func TestBool_NilFnFails(t *testing.T) {
	c := Bool("nilbool", nil)
	res := c.Run(t.Context())
	if res.OK {
		t.Fatal("expected FAIL for nil fn")
	}
	if res.Error != "nil check function" {
		t.Errorf("unexpected error: %q", res.Error)
	}
}

func TestBool_TrueIsOK(t *testing.T) {
	c := Bool("flag", func(_ context.Context) (bool, error) {
		return true, nil
	})
	res := c.Run(t.Context())
	if !res.OK {
		t.Fatal("expected OK")
	}
}

func TestBool_FalseIsFail(t *testing.T) {
	c := Bool("flag", func(_ context.Context) (bool, error) {
		return false, nil
	})
	res := c.Run(t.Context())
	if res.OK {
		t.Fatal("expected FAIL")
	}
	if !strings.Contains(res.Error, "false") {
		t.Errorf("expected error about false, got %q", res.Error)
	}
}

func TestBool_ErrorIsFail(t *testing.T) {
	c := Bool("flag", func(_ context.Context) (bool, error) {
		return false, errors.New("oops")
	})
	res := c.Run(t.Context())
	if res.OK {
		t.Fatal("expected FAIL")
	}
	if res.Error != "oops" {
		t.Errorf("unexpected error: %q", res.Error)
	}
}

func TestBool_PanicRecovery(t *testing.T) {
	c := Bool("panic-bool", func(_ context.Context) (bool, error) {
		panic("kaboom")
	})
	res := c.Run(t.Context())
	if res.OK {
		t.Fatal("expected FAIL")
	}
	if !strings.Contains(res.Error, "panic: kaboom") {
		t.Errorf("expected panic message, got %q", res.Error)
	}
}

func TestBool_EmptyLabelDefaultsToCustom(t *testing.T) {
	c := Bool("", func(_ context.Context) (bool, error) { return true, nil })
	if c.Name() != "custom" {
		t.Errorf("expected name 'custom', got %q", c.Name())
	}
}
