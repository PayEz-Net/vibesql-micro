package version

import (
	"fmt"
	"runtime"
)

// Version information
// These can be overridden at build time using ldflags
var (
	// Version is the semantic version of VibeSQL
	Version = "1.0.0"

	// GitCommit is the git commit hash
	GitCommit = "dev"

	// BuildDate is the date the binary was built
	BuildDate = "unknown"

	// GoVersion is the version of Go used to build the binary
	GoVersion = runtime.Version()
)

// Info contains all version information
type Info struct {
	Version   string
	GitCommit string
	BuildDate string
	GoVersion string
	OS        string
	Arch      string
}

// Get returns the version information
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf("VibeSQL %s (commit: %s, built: %s, go: %s, %s/%s)",
		i.Version,
		i.GitCommit,
		i.BuildDate,
		i.GoVersion,
		i.OS,
		i.Arch,
	)
}

// Short returns a short version string (version only)
func (i Info) Short() string {
	return i.Version
}

// Full returns a detailed multi-line version string
func (i Info) Full() string {
	return fmt.Sprintf(`VibeSQL Version Information:
  Version:    %s
  Git Commit: %s
  Build Date: %s
  Go Version: %s
  OS/Arch:    %s/%s`,
		i.Version,
		i.GitCommit,
		i.BuildDate,
		i.GoVersion,
		i.OS,
		i.Arch,
	)
}
