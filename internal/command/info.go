package command

import "github.com/spf13/cobra"

var infoCommand = cobra.Command{
	Use:   "info",
	Short: "Display information about an tool",
}

func init() {
	rootCommand.AddCommand(&infoCommand)
}
