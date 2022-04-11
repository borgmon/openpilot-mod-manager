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
	"github.com/borgmon/openpilot-mod-manager/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const OMMVersion = "v0.1"

var versionFlag = false

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "omm",
	Short: "Openpilot Mod Manager(OMM)",
	Long:  `Openpilot Mod Manager is a toolkit for using and developing mods.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PreRunE: func(cmd *cobra.Command, args []string) error { return loadParam() },
	RunE: func(cmd *cobra.Command, args []string) error {
		if versionFlag {
			fmt.Println(version.OMMVersion)
		}
		return cmd.Help()
	},
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

	if opPath, ok := os.LookupEnv("OPENPILOT_PATH"); ok {
		OPPath = opPath
	} else {
		OPPath = wd
	}
	if ommPath, ok := os.LookupEnv("OMM_PATH"); ok {
		OMMPath = ommPath
	} else {
		OMMPath = filepath.Join(home, config.CACHEPATH)
	}

	rootCmd.PersistentFlags().StringVarP(&OPPath, "openpilot", "o", OPPath, "openpilot path")
	rootCmd.PersistentFlags().StringVarP(&OMMPath, "omm", "m", OMMPath, "omm path")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "display extra info")
	rootCmd.Flags().BoolVar(&versionFlag, "version", false, "OMM version")
}

func loadParam() error {
	param.NewParam(
		filepath.Join(OPPath, config.CONFIG_FILE_NAME),
		OMMPath,
		OPPath,
		Verbose,
	)
	return file.GetFileHandler().NewFolder(OMMPath)
}

func load() error {
	err := loadParam()
	if err != nil {
		return err
	}
	_, err = config.LoadConfigHandler()
	if err != nil {
		e := errors.Unwrap(err)
		if _, ok := e.(*os.PathError); ok {
			return errors.New(fmt.Sprintf("No %v found in the directory %v.\nPlease init first:\nomm init", config.CONFIG_FILE_NAME, OPPath))
		}
		return err
	}

	return nil
}
