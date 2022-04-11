/*
Copyright © 2022 borgmon

*/
package cmd

import (
	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/installer"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// infoCmd represents the remove command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Info installed mods",
	Example: `omm info omm-base
omm info https://github.com/borgmon/omm-base
omm info /local/omm-base
`,
	PreRunE: func(cmd *cobra.Command, args []string) error { return load() },
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Invalid Args")
		}
		return common.LogIfErr(installer.GetInstaller().Info(args[0]))
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
