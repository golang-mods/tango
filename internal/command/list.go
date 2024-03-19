package command

import (
	"fmt"

	"github.com/golang-mods/tango"
	"github.com/golang-mods/tango/internal/process/manifest"
	"github.com/spf13/cobra"
)

var listCommand = cobra.Command{
	Use:   "list",
	Short: "List installed tools",
	RunE: createRun(func(command *cobra.Command, arguments []string, manager *tango.Manager) error {
		binaries, err := manager.List()
		if err != nil {
			return err
		}

		for _, binary := range binaries {
			fmt.Println(manifest.ToPackage(binary.Path, binary.Version.Module.Version))
		}

		return nil
	}),
}

func init() {
	rootCommand.AddCommand(&listCommand)
}
