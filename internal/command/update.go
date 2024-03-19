package command

import (
	"github.com/golang-mods/tango"
	"github.com/spf13/cobra"
)

var updateCommand = cobra.Command{
	Use:   "update [package]...",
	Short: "Update tools",
	Example: examplePrefix + " update\n" +
		examplePrefix + " update " + examplePath,
	RunE: createRun(func(command *cobra.Command, arguments []string, manager *tango.Manager) error {
		return manager.Update(arguments)
	}),
}

func init() {
	rootCommand.AddCommand(&updateCommand)
}
