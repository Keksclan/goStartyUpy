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

func TestLarsMessageInBanner(t *testing.T) {
	// Today is March 6, so isLarsBirthday() should be true and the message should appear.
	out := Render(Options{ServiceName: "test"}, BuildInfo{})
	if !isLarsBirthday() {
		t.Skip("not March 6 – skipping live banner check")
	}
	if got := out; !contains(got, "Lars hat Geburtstag") {
		t.Errorf("expected Lars birthday message in banner, got:\n%s", got)
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
