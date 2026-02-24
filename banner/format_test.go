package banner

import (
	"strings"
	"testing"
)

func TestAlignedColons(t *testing.T) {
	info := BuildInfo{
		Version:   "v1.0.0",
		BuildTime: "now",
		Commit:    "abc",
		Branch:    "main",
		Dirty:     "false",
	}
	opts := Options{
		ServiceName: "svc",
		Environment: "dev",
		Extra: map[string]string{
			"HTTP": ":8080",
			"gRPC": ":9090",
		},
	}

	out := Render(opts, info)

	// Every key-value line (indented with two spaces) must have its colon at
	// the same column, proving alignment is stable.
	colonCols := map[int]bool{}
	for _, line := range strings.Split(out, "\n") {
		if !strings.HasPrefix(line, "  ") {
			continue
		}
		idx := strings.Index(line, " : ")
		if idx < 0 {
			continue
		}
		colonCols[idx] = true
	}
	if len(colonCols) != 1 {
		t.Errorf("expected all colons at same column, got columns %v", colonCols)
	}
}

func TestExtra_SortedByKey(t *testing.T) {
	opts := Options{
		Extra: map[string]string{
			"Zulu":  "last",
			"Alpha": "first",
			"Mike":  "middle",
		},
	}
	out := Render(opts, BuildInfo{})

	idxA := strings.Index(out, "Alpha")
	idxM := strings.Index(out, "Mike")
	idxZ := strings.Index(out, "Zulu")

	if idxA < 0 || idxM < 0 || idxZ < 0 {
		t.Fatal("one or more Extra keys missing from output")
	}
	if !(idxA < idxM && idxM < idxZ) {
		t.Errorf("Extra keys not sorted: Alpha@%d, Mike@%d, Zulu@%d", idxA, idxM, idxZ)
	}
}

func TestExtra_EmptyMap(t *testing.T) {
	opts := Options{Extra: map[string]string{}}
	out := Render(opts, BuildInfo{})
	if out == "" {
		t.Error("output should not be empty with empty Extra map")
	}
}

func TestBuildKVs_ContainsStandardFields(t *testing.T) {
	opts := Options{
		ServiceName: "my-svc",
		Environment: "prod",
		Extra: map[string]string{
			"HTTP": ":80",
			"gRPC": ":443",
		},
	}
	info := BuildInfo{
		Version:   "v2.0.0",
		BuildTime: "2026-01-01",
		Commit:    "deadbeef",
		Branch:    "release",
		Dirty:     "true",
	}

	kvs := buildKVs(opts, info)
	keys := make([]string, len(kvs))
	for i, p := range kvs {
		keys[i] = p.Key
	}

	// Verify expected standard keys appear in order.
	expected := []string{"Service", "Environment", "Version",
		"BuildTime", "Commit", "Branch", "Dirty", "Go", "OS/Arch", "PID"}
	for i, want := range expected {
		if keys[i] != want {
			t.Errorf("key[%d] = %q, want %q", i, keys[i], want)
		}
	}

	// Extra keys (HTTP, gRPC) should appear after standard fields, sorted.
	extraKeys := keys[len(expected):]
	if len(extraKeys) < 2 || extraKeys[0] != "HTTP" || extraKeys[1] != "gRPC" {
		t.Errorf("extra keys = %v, want [HTTP gRPC]", extraKeys)
	}
}
