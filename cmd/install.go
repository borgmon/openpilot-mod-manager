/*
Copyright © 2022 borgmon

*/
package cmd

import (
	"fmt"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/borgmon/openpilot-mod-manager/mod"
	"github.com/borgmon/openpilot-mod-manager/source"
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
		return installMod(args[0], isForce)
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

func downloadMod(s source.Source, force bool) (*manifest.Manifest, error) {
	if force {
		return s.DownloadToCache()
	}
	name, err := s.GetName()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if mod, _ := ConfigHandler.FindMod(name); mod != nil {
		fmt.Println("This Mod already Exist")
		return nil, nil
	} else {
		return s.DownloadToCache()
	}
}

func installMod(path string, force bool) error {
	var s source.Source
	if common.IsUrl(path) {
		s = &source.GitSource{
			RemoteUrl:  path,
			GitHandler: GitHandler,
			CachePath:  CachePath,
		}

	} else {
		s = &source.LocalSource{
			LocalPath: path,
			CachePath: CachePath,
		}
	}

	man, err := downloadMod(s, force)
	if err != nil {
		return errors.WithStack(err)
	}

	err = ConfigHandler.AddMod(&mod.Mod{
		Name:    man.Name,
		Version: man.Version,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
