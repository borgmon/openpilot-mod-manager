/*
Copyright © 2022 borgmon

*/
package cmd

import (
	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/installer"
	"github.com/spf13/cobra"
)

// resetCmd represents the install command
var resetCmd = &cobra.Command{
	Use:     "reset",
	Short:   "Reset openpilot repo",
	Example: `omm reset`,
	PreRun:  func(cmd *cobra.Command, args []string) { populate() },
	RunE: func(cmd *cobra.Command, args []string) error {
		return common.LogIfErr(installer.GetInstaller().Reset())
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(resetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
