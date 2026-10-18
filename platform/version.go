package platform

import (
	"runtime/debug"
)

// VersionInfo is the build identity.
type VersionInfo struct {
	// Module is the main module path.
	Module string
	// Version is the VCS revision, else the module version, else "unknown".
	Version string
	// Commit is the VCS revision, when available.
	Commit string
	// Date is the VCS commit time, when available.
	Date string
	// Dirty reports whether the source tree had uncommitted changes.
	Dirty bool
}

// Version returns the build identity for the named module.
func Version(name string) *VersionInfo {
	info := VersionInfo{Version: "unknown"}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		info.Module = name
		return &info
	}

	info.Module = bi.Main.Path
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			info.Commit = s.Value
		case "vcs.time":
			info.Date = s.Value
		case "vcs.modified":
			info.Dirty = s.Value == "true"
		}
	}

	switch {
	case info.Commit != "":
		info.Version = info.Commit
	case bi.Main.Version != "" && bi.Main.Version != "(devel)":
		info.Version = bi.Main.Version
	}

	return &info
}

// String reports the module path and version.
func (info *VersionInfo) String() string {
	version := info.Version
	if info.Commit != "" {
		if info.Date != "" {
			version += " (committed " + info.Date + ")"
		}
		if info.Dirty {
			version += " [dirty]"
		}
	}
	return info.Module + " " + version
}
