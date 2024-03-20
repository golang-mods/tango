package command

import "github.com/golang-mods/tango/internal/constant"

const (
	shortDescription = "Tools manager for Go"
	highlight        = "^      ^^         ^^"
	longDescription  = shortDescription + "\n" + highlight

	examplePath            = "example.com/example/command"
	examplePathWithVersion = examplePath + "@^1.2.0"
	examplePrefix          = "  " + constant.ApplicationName
)
