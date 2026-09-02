package version

import (
	"runtime/debug"
	"strings"
)

// Value is set by release builds through -ldflags. Local tagged builds use VCS metadata.
var Value string

func String() string {
	if Value != "" {
		return Value
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.tag" && setting.Value != "" {
				return strings.TrimPrefix(setting.Value, "v")
			}
		}
	}
	return "0.0.0"
}
