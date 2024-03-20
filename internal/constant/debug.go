//go:build !release

package constant

import "runtime/debug"

func init() {
	info, ok := debug.ReadBuildInfo()
	if ok {
		ApplicationVersion = info.Main.Version
	}
}
