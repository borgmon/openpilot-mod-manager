package installer

import (
	"fmt"
	"path/filepath"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/config"
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/git"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/borgmon/openpilot-mod-manager/mod"
	"github.com/borgmon/openpilot-mod-manager/source"
	"github.com/pkg/errors"
)

type InstallerImpl struct {
	ConfigHandler config.ConfigHandler
	GitHandler    git.GitHandler
	OPPath        string
	CachePath     string
}

func NewInstaller(ConfigHandler config.ConfigHandler, GitHandler git.GitHandler, OPPath string, CachePath string) Installer {
	return &InstallerImpl{
		ConfigHandler: ConfigHandler,
		GitHandler:    GitHandler,
		OPPath:        OPPath,
		CachePath:     CachePath,
	}
}

func (installer *InstallerImpl) Apply() error {
	err := installer.GitHandler.CheckoutBranch(installer.OPPath, installer.ConfigHandler.GetConfig().OPVersion)
	if err != nil {
		return errors.WithStack(err)
	}
	err = installer.GitHandler.NewBranch(installer.OPPath, installer.GitHandler.GenerateBranchName())
	if err != nil {
		return errors.WithStack(err)
	}
	err = installer.ConfigHandler.ApplyMods()
	if err != nil {
		return errors.WithStack(err)
	}
	err = installer.GitHandler.CommitBranch(installer.OPPath, installer.GitHandler.GenerateBranchName())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (installer *InstallerImpl) Reset() error {
	err := installer.GitHandler.CheckoutBranch(installer.OPPath, installer.ConfigHandler.GetConfig().OPVersion)
	if err != nil {
		return errors.WithStack(err)
	}
	err = file.GetFileHandler().RemoveFolder(installer.CachePath)
	if err != nil {
		return errors.WithStack(err)
	}
	err = installer.ConfigHandler.CreateConfig()
	if err != nil {
		return errors.WithStack(err)
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
	if mod, _ := installer.ConfigHandler.FindMod(name); mod != nil {
		fmt.Println("This Mod is already exist")
		return nil, nil
	} else {
		return s.DownloadToCache()
	}
}

func (installer *InstallerImpl) Remove(name string) error {
	err := installer.ConfigHandler.RemoveMod(name)
	if err != nil {
		return errors.WithStack(err)
	}
	err = file.GetFileHandler().RemoveFolder(filepath.Join(installer.CachePath, name))
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
			RemoteUrl:  path,
			GitHandler: installer.GitHandler,
			CachePath:  installer.CachePath,
		}

	} else {
		s = &source.LocalSource{
			LocalPath: path,
			CachePath: installer.CachePath,
		}
	}

	man, err := installer.DownloadMod(s, force)
	if err != nil {
		return errors.WithStack(err)
	}
	if man == nil {
		return nil
	}

	err = installer.Apply()
	if err != nil {
		return errors.WithStack(err)
	}

	err = installer.ConfigHandler.AddMod(&mod.Mod{
		Name:    man.Name,
		Version: man.Version,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
