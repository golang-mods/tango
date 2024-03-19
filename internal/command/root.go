package command

import (
	"fmt"
	"os"

	"github.com/golang-mods/exerrors"
	"github.com/golang-mods/serrors"
	"github.com/golang-mods/tango"
	"github.com/golang-mods/tango/internal/constant"
	"github.com/golang-mods/tango/internal/log"
	"github.com/spf13/cobra"
)

var rootCommand = cobra.Command{
	Use:               constant.ApplicationName,
	Short:             shortDescription,
	Long:              longDescription,
	Version:           constant.ApplicationVersion,
	PersistentPreRunE: persistentPreRun,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var logger = log.NewLogger(false, log.ColorAuto)

func init() {
	currentDriectory, err := os.Getwd()
	fatalIfError(&rootCommand, err)

	flags := rootCommand.PersistentFlags()
	flags.BoolP("color", "c", false, "enable color mode")
	flags.BoolP("debug", "d", false, "enable debug mode")
	flags.StringP("dir", "C", currentDriectory, "change to directory")
}

func persistentPreRun(command *cobra.Command, arguments []string) error {
	flags := command.Flags()

	color := log.ColorAuto
	if flags.Changed("color") {
		if flag, err := flags.GetBool("color"); err != nil {
			return err
		} else if flag {
			color = log.ColorAlways
		} else {
			color = log.ColorNever
		}
	}

	debug, err := flags.GetBool("debug")
	if err != nil {
		return err
	}

	logger = log.NewLogger(debug, color)

	return nil
}

func Execute() {
	fatalIfError(rootCommand.ExecuteC())
}

func fatalIfError(command *cobra.Command, err error) {
	errs := exerrors.Flatten(err)
	if len(errs) == 0 {
		return
	}

	for _, err := range errs {
		logger.Error(err.Error(), serrors.Attributes(err)...)
	}
	fmt.Println()
	command.Usage()

	os.Exit(1)
}

func createRun(
	run func(command *cobra.Command, arguments []string, manager *tango.Manager) error,
) func(*cobra.Command, []string) error {
	return func(command *cobra.Command, arguments []string) error {
		flags := command.Flags()

		currentDirectory, err := flags.GetString("dir")
		if err != nil {
			return err
		}

		manager, err := tango.NewManager(currentDirectory, logger)
		if err != nil {
			return err
		}

		return run(command, arguments, manager)
	}
}
