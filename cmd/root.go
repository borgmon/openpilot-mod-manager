/*
Copyright © 2022 borgmon

*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/config"
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/git"
	"github.com/borgmon/openpilot-mod-manager/installer"
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

var ConfigHandler config.ConfigHandler
var GitHandler git.GitHandler
var Installer installer.Installer
var CachePath = ""
var ConfigPath = ""

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	wd, err := os.Getwd()
	common.PanicIfErr(err)
	home, err := os.UserHomeDir()
	common.PanicIfErr(err)
	rootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", filepath.Join(wd, config.CONFIG_FILE_NAME), "config file omm.yml")
	rootCmd.PersistentFlags().StringVarP(&CachePath, "cache", "a", filepath.Join(home, config.CACHEPATH), "cache dir")

	ConfigHandler = config.NewConfigHandler(ConfigPath, CachePath, "/home/kev/dev/omm-test")
	c, _ := ConfigHandler.LoadConfig()
	if c == nil {
		c, err = ConfigHandler.CreateConfig()
		common.PanicIfErr(err)
	}

	GitHandler = git.NewGitHandler(CachePath, c.OPVersion)

	err = file.GetFileHandler().NewFolder(CachePath)
	common.PanicIfErr(err)

	Installer = installer.NewInstaller(ConfigHandler, GitHandler, "/home/kev/dev/omm-test", CachePath)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
