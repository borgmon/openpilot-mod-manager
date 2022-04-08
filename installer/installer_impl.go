package installer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/config"
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/git"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/borgmon/openpilot-mod-manager/mod"
	"github.com/borgmon/openpilot-mod-manager/param"
	"github.com/borgmon/openpilot-mod-manager/source"
	"github.com/pkg/errors"
)

type InstallerImpl struct{}

var installerInstance Installer

func GetInstaller() Installer {
	if installerInstance != nil {
		return installerInstance
	}
	installerInstance = &InstallerImpl{}
	return installerInstance
}

func (installer *InstallerImpl) Apply() error {
	err := git.GetGitHandler().CheckoutBranch(param.PathStore.OPPath, config.GetConfigHandler().GetConfig().OPVersion)
	if err != nil {
		return errors.WithStack(err)
	}
	err = git.GetGitHandler().NewBranch(param.PathStore.OPPath, git.GetGitHandler().GenerateBranchName())
	if err != nil {
		return errors.WithStack(err)
	}
	err = config.GetConfigHandler().ApplyMods()
	if err != nil {
		return errors.WithStack(err)
	}
	err = git.GetGitHandler().AddBranch(param.PathStore.OPPath)
	if err != nil {
		return errors.WithStack(err)
	}
	err = git.GetGitHandler().CommitBranch(param.PathStore.OPPath, config.GetConfigHandler().BuildModList())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (installer *InstallerImpl) Reset() error {
	err := git.GetGitHandler().CheckoutBranch(param.PathStore.OPPath, config.GetConfigHandler().GetConfig().OPVersion)
	if err != nil {
		return errors.WithStack(err)
	}
	err = installer.RemoveAllOMMBranches()
	if err != nil {
		return errors.WithStack(err)
	}
	err = file.GetFileHandler().RemoveFolder(param.PathStore.OMMPath)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = config.GetConfigHandler().CreateConfig()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (installer *InstallerImpl) RemoveAllOMMBranches() error {
	str, err := git.GetGitHandler().ListBranch(param.PathStore.OPPath)
	if err != nil {
		return errors.WithStack(err)
	}
	branches := strings.Split(str, "\n")
	for _, b := range branches {
		if strings.Contains(b, "omm-") {
			err = git.GetGitHandler().RemoveBranch(param.PathStore.OPPath, b[2:])
			if err != nil {
				return errors.WithStack(err)
			}
		}

	}
	return nil
}

func (installer *InstallerImpl) DownloadMod(s source.Source, force bool) (*manifest.Manifest, error) {
	if force {
		return s.DownloadToCache()
	}
	name, err := s.GetName()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if mod, _ := config.GetConfigHandler().FindMod(name); mod != nil {
		fmt.Println("This Mod is already exist")
		return nil, nil
	} else {
		return s.DownloadToCache()
	}
}

func (installer *InstallerImpl) Remove(name string) error {
	err := config.GetConfigHandler().RemoveMod(name)
	if err != nil {
		return errors.WithStack(err)
	}
	err = file.GetFileHandler().RemoveFolder(filepath.Join(param.PathStore.OMMPath, name))
	if err != nil {
		return errors.WithStack(err)
	}
	err = installer.Apply()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (installer *InstallerImpl) Install(path string, force bool) error {
	var s source.Source
	if common.IsUrl(path) {
		s = &source.GitSource{
			RemoteUrl:     path,
			ConfigHandler: config.GetConfigHandler(),
		}

	} else {
		s = &source.LocalSource{
			LocalPath:     path,
			ConfigHandler: config.GetConfigHandler(),
		}
	}

	man, err := installer.DownloadMod(s, force)
	if err != nil {
		return errors.WithStack(err)
	}
	if man == nil {
		return nil
	}

	err = config.GetConfigHandler().AddMod(&mod.Mod{
		Name:    man.Name,
		Version: man.Version,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	err = installer.Apply()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (installer *InstallerImpl) List() error {
	_, err := fmt.Println(config.GetConfigHandler().BuildModList())
	return err
}
