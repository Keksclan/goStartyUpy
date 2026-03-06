package banner

import (
	"testing"
	"time"
)

func TestIsLarsBirthdayAt_March6(t *testing.T) {
	bd := time.Date(2026, time.March, 6, 12, 0, 0, 0, time.UTC)
	if !isLarsBirthdayAt(bd) {
		t.Error("expected true on March 6")
	}
}

func TestIsLarsBirthdayAt_OtherDay(t *testing.T) {
	other := time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC)
	if isLarsBirthdayAt(other) {
		t.Error("expected false on January 15")
	}
}

func TestLarsChanceConstant(t *testing.T) {
	if larsChance != 100 {
		t.Errorf("expected larsChance=100, got %d", larsChance)
	}
}

func TestIsLarsBirthday_Probabilistic(t *testing.T) {
	// On March 6 the function should return true roughly 1% of the time.
	// We call it many times and verify it fires at least once and not always.
	if !isLarsBirthdayAt(time.Now()) {
		t.Skip("not March 6 – skipping probabilistic check")
	}
	hits := 0
	const runs = 10_000
	for range runs {
		if isLarsBirthday() {
			hits++
		}
	}
	if hits == 0 {
		t.Error("expected at least one hit in 10 000 runs (p≈1%)")
	}
	if hits == runs {
		t.Error("got hit every single time – randomness broken")
	}
	// Loose sanity: expect between 10 and 500 hits (1% of 10k = 100).
	if hits < 10 || hits > 500 {
		t.Errorf("hits=%d out of %d – outside plausible range", hits, runs)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && searchString(s, sub)
}

func searchString(s, sub string) bool {
	for i := range len(s) - len(sub) + 1 {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
