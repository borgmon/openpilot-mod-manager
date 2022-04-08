/*
Copyright © 2022 borgmon

*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/borgmon/openpilot-mod-manager/config"
	"github.com/borgmon/openpilot-mod-manager/git"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "omm",
	Short: "Openpilot Mod Manager(OMM)",
	Long:  `Openpilot Mod Manager is a toolkit for using and developing mods.`,
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

var a string

var ConfigHandler config.ConfigHandler
var GitHandler git.GitHandler
var ConfigName = "omm.yml"
var CachePath = "~/.omm"
var OPPath = "/data/openpilot"
var ConfigPath = filepath.Join(OPPath, ConfigName)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", filepath.Join(OPPath, ConfigName), "config file omm.yml")
	rootCmd.PersistentFlags().StringVarP(&CachePath, "cache", "a", "~/.omm", "cache dir")

	ConfigHandler = &config.ConfigHandlerImpl{
		ConfigPath: ConfigPath,
		CachePath:  CachePath,
		OPPath:     OPPath,
	}

	GitHandler = &git.GitHandlerImpl{
		CachePath: CachePath,
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
