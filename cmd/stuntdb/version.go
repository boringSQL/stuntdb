package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

var version = "dev"

func build() string {
	commit, goVersion := "none", runtime.Version()
	if info, ok := debug.ReadBuildInfo(); ok {
		if version == "dev" {
			if v := info.Main.Version; v != "" && v != "(devel)" && !strings.HasSuffix(v, "+dirty") {
				version = v
			}
		}
		if info.GoVersion != "" {
			goVersion = info.GoVersion
		}
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" && setting.Value != "" {
				commit = setting.Value[:min(7, len(setting.Value))]
			}
		}
	}
	return fmt.Sprintf("stuntdb %s (commit %s, %s)", version, commit, goVersion)
}
