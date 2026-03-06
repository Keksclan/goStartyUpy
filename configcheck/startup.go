package configcheck

import (
	"fmt"
	"os"
)

// Options controls whether configuration validation runs during startup.
type Options struct {
	// Enabled activates config validation. When false (default), the
	// [RunStartupCheck] function is a no-op.
	Enabled bool
	// Config is the configuration struct (or pointer to struct) to
	// validate. It is typically the value returned by goConfy's Load.
	Config any
}

// RunStartupCheck validates the configuration according to opts. It
// returns nil when validation is disabled or the configuration is valid.
// On failure it returns a *[ValidationError] whose Error() method
// produces a human-readable diagnostic block.
//
// When validation passes, a green checkmark line is printed to stdout
// so the operator gets immediate visual feedback (similar to the
// [OK] tags produced by the checks package).
func RunStartupCheck(opts Options) error {
	if !opts.Enabled {
		return nil
	}
	if opts.Config == nil {
		return &ValidationError{
			Errors: []string{"configuration is nil"},
		}
	}
	ve := Validate(opts.Config)
	if ve != nil {
		return ve
	}
	printSuccess(os.Stdout)
	return nil
}

// MustPassStartupCheck is like [RunStartupCheck] but calls the provided
// fatalf function (typically log.Fatalf) when validation fails. It is a
// convenience wrapper for use in main functions.
func MustPassStartupCheck(opts Options, fatalf func(string, ...any)) {
	if err := RunStartupCheck(opts); err != nil {
		fatalf("startup config check failed:\n%s", err)
	}
}

// ANSI escape sequences for the success indicator.
const (
	esc            = "\033["
	greenBold      = esc + "1;92m"
	reset          = esc + "0m"
	checkMark      = "✔"
	successMessage = "Config Check"
)

// printSuccess writes a green checkmark line to w.
func printSuccess(w *os.File) {
	if fileSupportsColor(w) {
		fmt.Fprintf(w, "  %s%s%s  %s%s%s\n", greenBold, checkMark, reset, greenBold, successMessage, reset)
	} else {
		fmt.Fprintf(w, "  [OK]  %s\n", successMessage)
	}
}

// fileSupportsColor reports whether the given file is a terminal that
// likely supports ANSI colors. On non-TTY outputs (pipes, files) it
// returns false so we fall back to plain text.
func fileSupportsColor(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// FormatValidationError returns a formatted string suitable for
// terminal output describing the validation problems.
func FormatValidationError(err error) string {
	if err == nil {
		return ""
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		return fmt.Sprintf("Config validation error: %v\n", err)
	}
	return ve.Error()
}
