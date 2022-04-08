/*
Copyright © 2022 borgmon

*/
package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove [mod name]",
	Short:   "Remove a mod",
	Example: `omm remove my-awesome-mod`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Invalid Args")
		}
		return LogIfErr(Installer.Remove(args[0]))
	},
	Aliases: []string{"i"},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}