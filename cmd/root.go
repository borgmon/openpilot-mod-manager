/*
Copyright © 2022 borgmon

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/config"
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/param"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "omm",
	Short:        "Openpilot Mod Manager(OMM)",
	Long:         `Openpilot Mod Manager is a toolkit for using and developing mods.`,
	SilenceUsage: true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var OMMPath string
var OPPath string
var Verbose bool

func init() {
	wd, err := os.Getwd()
	common.PanicIfErr(err)
	home, err := os.UserHomeDir()
	common.PanicIfErr(err)

	rootCmd.PersistentFlags().StringVarP(&OPPath, "openpilot", "o", wd, "openpilot path")
	rootCmd.PersistentFlags().StringVarP(&OMMPath, "omm", "m", filepath.Join(home, config.CACHEPATH), "omm path")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "display extra info")
}

func loadParam() error {
	param.NewParam(
		filepath.Join(OPPath, config.CONFIG_FILE_NAME),
		OMMPath,
		OPPath,
		Verbose,
	)
	_, err := config.LoadConfigHandler()
	if err != nil {
		e := errors.Unwrap(err)
		if _, ok := e.(*os.PathError); ok {
			return errors.New(fmt.Sprintf("No %v found in the directory %v.\nPlease init first:\nomm init\n", config.CONFIG_FILE_NAME, OPPath))
		}
		return err
	}

	return file.GetFileHandler().NewFolder(OMMPath)
}
