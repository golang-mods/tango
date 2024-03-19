//go:build !release

package constant

import (
	"path"
	"runtime/debug"
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if ok {
		ApplicationName = path.Base(info.Main.Path)
		ApplicationVersion = info.Main.Version
	}
}
