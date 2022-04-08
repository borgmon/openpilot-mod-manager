/*
Copyright © 2022 borgmon

*/
package cmd

import (
	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [url/local path]",
	Short: "Install a mod",
	Example: `omm install https://github.com/borgmon/omm-no-disengage_on_gas
omm install /home/usr/my-awesome-mod`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Invalid Args")
		}
		isForce, err := cmd.Flags().GetBool("force")
		if err != nil {
			return errors.WithStack(err)
		}
		return common.LogIfErr(Installer.Install(args[0], isForce))
	},
	Aliases: []string{"i"},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	installCmd.Flags().BoolP("force", "f", false, "force install a mod")
}
