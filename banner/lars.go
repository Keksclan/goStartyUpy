package banner

import (
	"math/rand/v2"
	"time"
)

// larsMessage is the birthday greeting for Lars. 🎂
const larsMessage = "🎉 Lars hat Geburtstag! 🎂"

// larsChance is the 1-in-N probability of showing the birthday message.
const larsChance = 100

// isLarsBirthday returns true on March 6th with a 1:100 chance.
func isLarsBirthday() bool {
	now := time.Now()
	if now.Month() != time.March || now.Day() != 6 {
		return false
	}
	return rand.IntN(larsChance) == 0
}

// isLarsBirthdayAt returns true when t falls on March 6th (testable variant,
// ignores the random chance).
func isLarsBirthdayAt(t time.Time) bool {
	return t.Month() == time.March && t.Day() == 6
}
