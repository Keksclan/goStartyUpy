package checks

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestGroup_AllOK(t *testing.T) {
	g := NewGroup("infra", GroupOptions{}, []Check{
		New("a", func(_ context.Context) error { return nil }),
		New("b", func(_ context.Context) error { return nil }),
	}...)
	res := g.Run(t.Context())
	if !res.OK {
		t.Fatalf("expected OK, got error: %s", res.Error)
	}
	if res.Name != "infra" {
		t.Errorf("expected name 'infra', got %q", res.Name)
	}
}

func TestGroup_FailingChildMakesGroupFail(t *testing.T) {
	g := NewGroup("deps", GroupOptions{},
		New("ok-child", func(_ context.Context) error { return nil }),
		New("bad-child", func(_ context.Context) error { return errors.New("down") }),
	)
	res := g.Run(t.Context())
	if res.OK {
		t.Fatal("expected FAIL when a child fails")
	}
	if !strings.Contains(res.Error, "bad-child") {
		t.Errorf("error should mention failing child name, got %q", res.Error)
	}
	if !strings.Contains(res.Error, "down") {
		t.Errorf("error should mention child error, got %q", res.Error)
	}
	if !strings.Contains(res.Error, "1 failing") {
		t.Errorf("error should contain count, got %q", res.Error)
	}
}

func TestGroup_MultipleFailures(t *testing.T) {
	g := NewGroup("multi", GroupOptions{},
		New("fail1", func(_ context.Context) error { return errors.New("e1") }),
		New("ok1", func(_ context.Context) error { return nil }),
		New("fail2", func(_ context.Context) error { return errors.New("e2") }),
	)
	res := g.Run(t.Context())
	if res.OK {
		t.Fatal("expected FAIL")
	}
	if !strings.Contains(res.Error, "2 failing") {
		t.Errorf("expected '2 failing' in error, got %q", res.Error)
	}
	if !strings.Contains(res.Error, "fail1") || !strings.Contains(res.Error, "fail2") {
		t.Errorf("error should list both failing children, got %q", res.Error)
	}
}

func TestGroup_EmptyChecksIsOK(t *testing.T) {
	g := NewGroup("empty", GroupOptions{})
	res := g.Run(t.Context())
	if !res.OK {
		t.Fatal("expected OK for empty group")
	}
}

func TestGroup_EmptyLabelDefaultsToGroup(t *testing.T) {
	g := NewGroup("", GroupOptions{},
		New("x", func(_ context.Context) error { return nil }),
	)
	if g.Name() != "group" {
		t.Errorf("expected name 'group', got %q", g.Name())
	}
}

func TestGroup_WithParallelOption(t *testing.T) {
	g := NewGroup("parallel-deps", GroupOptions{Parallel: true, TimeoutPerCheck: time.Second},
		New("a", func(_ context.Context) error { return nil }),
		New("b", func(_ context.Context) error { return nil }),
	)
	res := g.Run(t.Context())
	if !res.OK {
		t.Fatalf("expected OK, got error: %s", res.Error)
	}
}

func TestDefaultRunner(t *testing.T) {
	r := DefaultRunner()
	if r.TimeoutPerCheck != 2*time.Second {
		t.Errorf("expected 2s timeout, got %v", r.TimeoutPerCheck)
	}
	if !r.Parallel {
		t.Error("expected Parallel to be true")
	}
}
