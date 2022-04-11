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

// initCmd represents the remove command
var initCmd = &cobra.Command{
	Use:     "init",
	Short:   fmt.Sprintf("Init this directory and generate %v file", config.CONFIG_FILE_NAME),
	Example: `omm init`,
	PreRunE: func(cmd *cobra.Command, args []string) error { return loadParam() },
	RunE: func(cmd *cobra.Command, args []string) error {
		return common.LogIfErr(installer.GetInstaller().Init(OPPath))
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
