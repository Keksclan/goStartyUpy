package banner

import "time"

// larsMessage is the birthday greeting for Lars. 🎂
const larsMessage = "🎉 Lars hat Geburtstag! 🎂"

// isLarsBirthday returns true on March 6th.
func isLarsBirthday() bool {
	now := time.Now()
	return now.Month() == time.March && now.Day() == 6
}

// isLarsBirthdayAt returns true when t falls on March 6th (testable variant).
func isLarsBirthdayAt(t time.Time) bool {
	return t.Month() == time.March && t.Day() == 6
}
