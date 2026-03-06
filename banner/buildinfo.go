// Package banner provides startup banner rendering with build metadata and
// runtime information. It is designed to be imported by services that want a
// consistent, informative startup message.
package banner

// Build metadata variables – override at link time with -ldflags:
//
//	go build -ldflags "-X 'github.com/keksclan/goStartyUpy/banner.Version=v1.2.3' ..."
var (
	// Version is the semantic version of the build.
	Version = "dev"
	// BuildTime is the UTC timestamp of when the binary was built.
	BuildTime = "unknown"
	// Commit is the short git commit hash of the source tree.
	Commit = "unknown"
	// Branch is the git branch the binary was built from.
	Branch = "unknown"
	// Dirty indicates whether the working tree had uncommitted changes ("true"/"false").
	Dirty = "false"
)

// BuildInfo holds the build metadata for the running binary.
// It is created via [CurrentBuildInfo] which snapshots only the
// package-level link-time variables (not runtime values).
type BuildInfo struct {
	// Version is the semantic version or git-describe output (e.g. "v1.2.3").
	// Defaults to "dev" when not set via -ldflags.
	Version string
	// BuildTime is the build timestamp in RFC 3339 format as provided via
	// -ldflags. Defaults to "unknown" when not set.
	BuildTime string
	// Commit is the short git commit hash of the source tree as provided
	// by the package-level variable. Defaults to "unknown" when not set
	// via -ldflags.
	Commit string
	// Branch is the git branch the binary was built from.
	// Defaults to "unknown" when not set via -ldflags.
	Branch string
	// Dirty is "true" when the working tree had uncommitted changes
	// at build time, "false" otherwise.
	Dirty string
}

// CurrentBuildInfo returns a [BuildInfo] snapshot populated exclusively from
// the package-level link-time variables. It does not capture any runtime
// values.
func CurrentBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		BuildTime: BuildTime,
		Commit:    Commit,
		Branch:    Branch,
		Dirty:     Dirty,
	}
}
