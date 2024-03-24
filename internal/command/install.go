package command

import (
	"github.com/golang-mods/tango"
	"github.com/golang-mods/tango/internal/constant"
	"github.com/spf13/cobra"
)

var installCommand = cobra.Command{
	Use:     "install [package]...",
	Aliases: []string{"i"},
	Short:   "Install tools",
	Example: examplePrefix + " install\n" +
		examplePrefix + " install " + examplePath + "\n" +
		examplePrefix + " install " + examplePathWithVersion + "\n\n" +
		`Version:
  ~1.2.0      fixed major version (default)
  ^1.2.0      fixed minor version
  *           latest version
  latest      same *
  reference   refer to ` + constant.GoModFile,
	RunE: createRun(func(command *cobra.Command, arguments []string, manager *tango.Manager) error {
		if len(arguments) == 0 {
			return manager.InstallAll()
		}

		flags := command.Flags()
		env, err := flags.GetStringArray("env")
		if err != nil {
			return err
		}

		return manager.Install(arguments, tango.InstallOptionEnv(env))
	}),
}

func init() {
	flags := installCommand.Flags()
	flags.StringArrayP("env", "e", nil, "set env")

	rootCommand.AddCommand(&installCommand)
}
