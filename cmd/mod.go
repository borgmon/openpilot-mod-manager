/*
Copyright © 2022 borgmon

*/
package cmd

import (
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// modCmd represents the mod command
var modCmd = &cobra.Command{
	Use:   "mod [subcommand]",
	Short: "Tools to develop mods",
	Example: `omm mod init
omm mod new`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a manifest.yml",
	Long: `Example:
omm mod init`,
	Example: "example",
	PreRun:  func(cmd *cobra.Command, args []string) { populate() },
	RunE: func(cmd *cobra.Command, args []string) error {
		return manifest.GetManifestHandler().Init(ConfigHandler.GetPaths().OPPath)
	},
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new patch",
	Long: `Example:
omm mod new selfdrive/common/params.cc`,
	Example: "example",
	PreRun:  func(cmd *cobra.Command, args []string) { populate() },
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Invalid Args")
		}
		return file.GetFileHandler().NewFileRecursively(args[0])
	},
}

func init() {
	modCmd.AddCommand(initCmd)
	modCmd.AddCommand(newCmd)
	rootCmd.AddCommand(modCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
