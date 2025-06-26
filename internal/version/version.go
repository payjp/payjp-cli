package version

import (
	"fmt"
	"runtime"
)

// These variables are set via ldflags during build
var (
	// Version is the semantic version of the CLI
	Version = "dev"

	// Commit is the git commit hash
	Commit = "none"

	// Date is the build date
	Date = "unknown"
)

// GetVersion returns the version string
func GetVersion() string {
	return Version
}

// GetFullVersion returns the full version string with build info
func GetFullVersion() string {
	return fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, Date)
}

// GetUserAgent returns the User-Agent string for HTTP requests
func GetUserAgent() string {
	return fmt.Sprintf("payjp-cli/%s go/%s", Version, runtime.Version())
}
