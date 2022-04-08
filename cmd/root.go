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

var CachePath = ""
var ConfigFilePath = ""

var ConfigHandler config.ConfigHandler
var Installer installer.Installer

func init() {
	wd, err := os.Getwd()
	common.PanicIfErr(err)
	home, err := os.UserHomeDir()
	common.PanicIfErr(err)

	rootCmd.PersistentFlags().StringVarP(&ConfigFilePath, "config", "c", filepath.Join(wd, config.CONFIG_FILE_NAME), "config file omm.yml")
	rootCmd.PersistentFlags().StringVarP(&CachePath, "cache", "a", filepath.Join(home, config.CACHEPATH), "cache dir")
}

func populate() {
	ConfigHandler = config.NewConfigHandler(&config.Paths{
		ConfigPath: ConfigFilePath,
		CachePath:  CachePath,
		OPPath:     filepath.Dir(ConfigFilePath)})
	c, _ := ConfigHandler.LoadConfig()
	if c == nil {
		err := ConfigHandler.SaveConfig()
		common.PanicIfErr(err)
	}

	err := file.GetFileHandler().NewFolder(CachePath)
	common.PanicIfErr(err)

	Installer = installer.NewInstaller(ConfigHandler)

}
