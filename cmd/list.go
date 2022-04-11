/*
Copyright © 2022 borgmon

*/
package cmd

import (
	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/installer"
	"github.com/spf13/cobra"
)

// listCmd represents the remove command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List installed mods",
	Example: `omm list`,
	PreRunE: func(cmd *cobra.Command, args []string) error { return loadParam() },
	RunE: func(cmd *cobra.Command, args []string) error {
		return common.LogIfErr(installer.GetInstaller().List())
	},
	Aliases:      []string{"l"},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
