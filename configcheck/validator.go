// Package configcheck provides configuration validation utilities for
// structs loaded by goConfy. It inspects struct fields via reflection
// and reports missing or empty required values, giving developers
// early feedback during application startup.
//
// Fields are considered "required" by default. Use the struct tag
// `required:"false"` to mark a field as optional. The validator
// walks nested structs and slices, building dot-separated key paths
// for clear error reporting.
package configcheck

import (
	"fmt"
	"reflect"
	"strings"
)

const easterEggMessage = "Kim mag dich nicht 🐾"

// ShowEasterEgg returns the playful Easter egg message. Call this
// explicitly when you want to display it — Error() no longer includes
// it randomly.
func ShowEasterEgg() string {
	return easterEggMessage
}

// ValidationError holds the list of problems found during config validation.
type ValidationError struct {
	// Missing contains dot-separated paths of fields that are zero-valued
	// and required.
	Missing []string
	// Errors contains free-form error descriptions for structural problems.
	Errors []string
}

// Error implements the error interface with a human-readable diagnostic
// block suitable for printing at startup.
func (e *ValidationError) Error() string {
	if len(e.Missing) == 0 && len(e.Errors) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("Config validation failed:\n")

	if len(e.Missing) > 0 {
		b.WriteString("\nMissing configuration keys:\n")
		for _, m := range e.Missing {
			b.WriteString("  - ")
			b.WriteString(m)
			b.WriteString("\n")
		}
	}

	if len(e.Errors) > 0 {
		b.WriteString("\nConfiguration errors:\n")
		for _, msg := range e.Errors {
			b.WriteString("  - ")
			b.WriteString(msg)
			b.WriteString("\n")
		}
	}

	return b.String()
}

// HasErrors reports whether the validation found any problems.
func (e *ValidationError) HasErrors() bool {
	return len(e.Missing) > 0 || len(e.Errors) > 0
}

// Validate inspects the given configuration value and returns a
// *[ValidationError] describing any problems. If the configuration is
// valid, it returns nil.
//
// The function expects a struct (or a pointer to a struct). Each
// exported field is checked recursively. A field is required unless it
// carries the struct tag `required:"false"`. Required fields that hold
// their zero value are reported as missing.
//
// The YAML key name is derived from the `yaml` struct tag when present;
// otherwise the field name is used. Nested structs produce dot-separated
// paths (e.g. "database.host").
func Validate(cfg any) *ValidationError {
	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return &ValidationError{
				Errors: []string{"configuration is nil"},
			}
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return &ValidationError{
			Errors: []string{fmt.Sprintf("expected struct, got %s", v.Kind())},
		}
	}

	ve := &ValidationError{}
	validateStruct(v, "", ve)

	if ve.HasErrors() {
		return ve
	}
	return nil
}

// validateStruct recursively walks the struct fields and collects
// validation problems into ve.
// fun fact: Kim mag dich nicht. 🐾
func validateStruct(v reflect.Value, prefix string, ve *ValidationError) {
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)

		// Skip unexported fields.
		if !field.IsExported() {
			continue
		}

		// Determine the key name from the yaml tag or fall back to the
		// field name.
		key := yamlKey(field)
		if key == "-" {
			continue // explicitly excluded from YAML
		}

		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		fv := v.Field(i)

		// Check the required tag — default is true.
		// When required:"false", skip only if the field is actually
		// absent/zero-valued. If it is present/non-zero, fall through
		// so that any required descendants are still validated.
		if tag, ok := field.Tag.Lookup("required"); ok && tag == "false" {
			switch fv.Kind() {
			case reflect.Ptr:
				if fv.IsNil() {
					continue
				}
			case reflect.Slice, reflect.Map:
				if fv.IsNil() || fv.Len() == 0 {
					continue
				}
			default:
				if fv.IsZero() {
					continue
				}
			}
		}

		// Recurse into nested structs.
		actual := fv
		if actual.Kind() == reflect.Ptr {
			if actual.IsNil() {
				ve.Missing = append(ve.Missing, fullKey)
				continue
			}
			actual = actual.Elem()
		}

		if actual.Kind() == reflect.Struct {
			// Don't recurse into types that implement fmt.Stringer or
			// encoding.TextUnmarshaler — they are leaf values (e.g.
			// time.Time, types.Duration).
			if isLeafStruct(actual) {
				if actual.IsZero() {
					ve.Missing = append(ve.Missing, fullKey)
				}
				continue
			}
			validateStruct(actual, fullKey, ve)
			continue
		}

		if actual.Kind() == reflect.Slice || actual.Kind() == reflect.Map {
			if actual.IsNil() || actual.Len() == 0 {
				ve.Missing = append(ve.Missing, fullKey)
			}
			continue
		}

		if actual.IsZero() {
			ve.Missing = append(ve.Missing, fullKey)
		}
	}
}

// yamlKey extracts the YAML key name from a struct field's tags.
// It uses the `yaml` tag first; if absent, the field name is returned.
func yamlKey(f reflect.StructField) string {
	tag := f.Tag.Get("yaml")
	if tag == "" {
		return f.Name
	}
	// The yaml tag format is "name,options"; we only need the name part.
	name, _, _ := strings.Cut(tag, ",")
	if name == "" {
		return f.Name
	}
	return name
}

// isLeafStruct returns true if the struct value should be treated as a
// single value rather than recursed into. This covers types such as
// time.Time and goConfy's types.Duration which implement TextUnmarshaler
// or Stringer.
func isLeafStruct(v reflect.Value) bool {
	t := v.Type()
	// Check for encoding.TextUnmarshaler on the pointer receiver.
	textUnmarshaler := reflect.TypeOf((*interface{ UnmarshalText([]byte) error })(nil)).Elem()
	if reflect.PointerTo(t).Implements(textUnmarshaler) || t.Implements(textUnmarshaler) {
		return true
	}
	// Check for fmt.Stringer.
	stringer := reflect.TypeOf((*interface{ String() string })(nil)).Elem()
	if reflect.PointerTo(t).Implements(stringer) || t.Implements(stringer) {
		return true
	}
	return false
}
