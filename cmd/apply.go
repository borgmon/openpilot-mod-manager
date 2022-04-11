/*
Copyright © 2022 borgmon

*/
package cmd

import (
	"fmt"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/config"
	"github.com/borgmon/openpilot-mod-manager/installer"
	"github.com/spf13/cobra"
)

// applyCmd represents the remove command
var applyCmd = &cobra.Command{
	Use:     "apply",
	Short:   fmt.Sprintf("Apply current %v", config.CONFIG_FILE_NAME),
	Example: `omm apply`,
	PreRunE: func(cmd *cobra.Command, args []string) error { return loadParam() },
	RunE: func(cmd *cobra.Command, args []string) error {
		return common.LogIfErr(installer.GetInstaller().Apply())
	},
	Aliases:      []string{"a"},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
