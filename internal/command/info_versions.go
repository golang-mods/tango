package command

import (
	"fmt"

	"github.com/golang-mods/tango"
	"github.com/spf13/cobra"
)

var infoVersionsCommand = cobra.Command{
	Use:     "versions package",
	Short:   "Display versions about an tool",
	Example: examplePrefix + "info versions " + examplePath,
	Args:    cobra.ExactArgs(1),
	RunE: createRun(func(command *cobra.Command, arguments []string, manager *tango.Manager) error {
		versions, err := manager.InfoVersions(arguments[0])
		if err != nil {
			return err
		}

		for _, version := range versions {
			fmt.Println(version)
		}

		return nil
	}),
}

func init() {
	infoCommand.AddCommand(&infoVersionsCommand)
}
