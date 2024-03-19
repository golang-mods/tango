package command

import (
	"github.com/golang-mods/tango"
	"github.com/golang-mods/tango/internal/constant"
	"github.com/spf13/cobra"
)

var initCommand = cobra.Command{
	Use:   "init",
	Short: "Create " + constant.ManifestFileName + " file",
	RunE: createRun(func(command *cobra.Command, arguments []string, manager *tango.Manager) error {
		return manager.Init()
	}),
}

func init() {
	rootCommand.AddCommand(&initCommand)
}
