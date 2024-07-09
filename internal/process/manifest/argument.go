package manifest

import (
	"errors"
	"strings"

	"github.com/golang-mods/tango/internal/gocommand"
	"github.com/samber/lo"
)

func ParseArgument(command *gocommand.Memorized, text string) (ManifestTool, error) {
	path := text
	versionText := ""
	if index := strings.LastIndex(text, "@"); index != -1 {
		path = text[:index]
		versionText = text[index+1:]
	}

	switch versionText {
	case "latest":
		versionText = "*"
	case "":
		versionText = "^"
		fallthrough
	case "^":
		fallthrough
	case "~":
		versions, err := command.MemorizedVersions(path)
		if err != nil {
			return ManifestTool{}, err
		}

		version, ok := lo.Last(versions.Versions)
		if !ok {
			return ManifestTool{}, errors.New("versions is empty")
		}

		versionText += version.String()
	}

	var tool ManifestTool
	if err := tool.Version.UnmarshalText([]byte(versionText)); err != nil {
		return ManifestTool{}, err
	}
	tool.Path = path

	return tool, nil
}
