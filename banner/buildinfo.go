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
type BuildInfo struct {
	Version   string
	BuildTime string
	Commit    string
	Branch    string
	Dirty     string
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
