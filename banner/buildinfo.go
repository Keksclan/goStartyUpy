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
	// Commit is the full git commit hash of the source tree.
	Commit = "unknown"
	// Branch is the git branch the binary was built from.
	Branch = "unknown"
	// Dirty indicates whether the working tree had uncommitted changes ("true"/"false").
	Dirty = "false"
)

// BuildInfo holds the build metadata for the running binary.
// It is created via [CurrentBuildInfo] which reads the link-time variables
// and runtime values into a single snapshot.
type BuildInfo struct {
	// Version is the semantic version or git-describe output (e.g. "v1.2.3").
	// Defaults to "dev" when not set via -ldflags.
	Version string
	// BuildTime is the UTC build timestamp in RFC 3339 format.
	// Defaults to "unknown" when not set via -ldflags.
	BuildTime string
	// Commit is the short git commit hash of the source tree.
	// Defaults to "unknown" when not set via -ldflags.
	Commit string
	// Branch is the git branch the binary was built from.
	// Defaults to "unknown" when not set via -ldflags.
	Branch string
	// Dirty is "true" when the working tree had uncommitted changes
	// at build time, "false" otherwise.
	Dirty string
}

// CurrentBuildInfo returns a BuildInfo snapshot populated from the
// link-time variables.
func CurrentBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		BuildTime: BuildTime,
		Commit:    Commit,
		Branch:    Branch,
		Dirty:     Dirty,
	}
}
