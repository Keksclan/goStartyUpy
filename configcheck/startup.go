package configcheck

import "fmt"

// Options controls whether configuration validation runs during startup.
type Options struct {
	// Enabled activates config validation. When false (default), the
	// [RunStartupCheck] function is a no-op.
	Enabled bool
	// Config is the configuration struct (or pointer to struct) to
	// validate. It is typically the value returned by goConfy's Load.
	Config any
	// Color enables ANSI color codes in the success message. When
	// false the returned string uses a plain-text "[OK]" tag instead.
	Color bool
}

// RunStartupCheck validates the configuration according to opts. It
// returns an empty string and nil error when validation is disabled.
// When the configuration is valid it returns a human-readable success
// message (colored or plain depending on [Options.Color]) and a nil
// error. On failure it returns an empty string and a *[ValidationError]
// whose Error() method produces a human-readable diagnostic block.
//
// The caller decides where (and whether) to print the success message,
// consistent with [banner.Render] which also returns a string.
func RunStartupCheck(opts Options) (string, error) {
	if !opts.Enabled {
		return "", nil
	}
	if opts.Config == nil {
		return "", &ValidationError{
			Errors: []string{"configuration is nil"},
		}
	}
	ve := Validate(opts.Config)
	if ve != nil {
		return "", ve
	}
	return formatSuccess(opts.Color), nil
}

// MustPassStartupCheck is like [RunStartupCheck] but calls the provided
// fatalf function (typically log.Fatalf) when validation fails. On
// success the message is printed to stdout for operator feedback. It is
// a convenience wrapper for use in main functions.
func MustPassStartupCheck(opts Options, fatalf func(string, ...any)) {
	msg, err := RunStartupCheck(opts)
	if err != nil {
		fatalf("startup config check failed:\n%s", err)
		return
	}
	if msg != "" {
		fmt.Print(msg)
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

// formatSuccess returns a success indicator line. When color is true
// the output contains ANSI escape sequences; otherwise plain text.
func formatSuccess(color bool) string {
	if color {
		return fmt.Sprintf("  %s%s%s  %s%s%s\n", greenBold, checkMark, reset, greenBold, successMessage, reset)
	}
	return fmt.Sprintf("  [OK]  %s\n", successMessage)
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
