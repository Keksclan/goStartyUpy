// Package version exposes the release version of the goStartyUpy module.
//
// ModuleVersion tracks the library's own semantic version and is updated
// with every tagged release. It is independent from the build-time
// metadata injected via ldflags into banner.Version, which represents the
// version of the service binary being built, not the library itself.
package version

// ModuleVersion is the semantic version of the goStartyUpy module.
// It must match the latest release entry in CHANGELOG.md and the
// corresponding git tag (e.g. v0.1.0).
//
// This value is set at source level and does not change at link time.
// For the service build version injected via ldflags, see banner.Version.
const ModuleVersion = "0.1.0"
