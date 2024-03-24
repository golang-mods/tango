package command

import (
	"github.com/golang-mods/tango"
	"github.com/spf13/cobra"
)

var uninstallCommand = cobra.Command{
	Use:     "uninstall package [package]...",
	Aliases: []string{"un"},
	Short:   "Uninstall tools",
	Example: examplePrefix + " uninstall " + examplePath,
	RunE: createRun(func(command *cobra.Command, arguments []string, manager *tango.Manager) error {
		return manager.Uninstall(arguments)
	}),
}

func init() {
	rootCommand.AddCommand(&uninstallCommand)
}
