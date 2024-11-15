package lib

import (
	"fmt"
	"runtime"
)

// Version holds the version information
type Version struct {
	Major      int
	Minor      int
	Patch      int
	CommitHash string
	BuildDate  string
}

// These variables will be set at build time using ldflags
var (
	commitHash string
	buildDate  string
)

// CurrentVersion holds the current version of the application
var CurrentVersion = Version{
	Major:      0,
	Minor:      1,
	Patch:      3,
	CommitHash: commitHash,
	BuildDate:  buildDate,
}

// String returns a string representation of the version
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// FullString returns a detailed string representation of the version
func (v Version) FullString() string {
	return fmt.Sprintf("Version: %s\nCommit: %s\nBuild Date: %s\nGo Version: %s\nOS/Arch: %s/%s",
		v.String(), v.CommitHash, v.BuildDate, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
