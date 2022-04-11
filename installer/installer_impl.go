package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/borgmon/openpilot-mod-manager/cache"
	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/config"
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/git"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/borgmon/openpilot-mod-manager/mod"
	"github.com/borgmon/openpilot-mod-manager/param"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
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
	fmt.Println("Applying...")
	err := git.GetGitHandler().CheckoutBranch(param.PathStore.OPPath, config.GetConfigHandler().GetConfig().OPVersion)
	if err != nil {
		return err
	}
	err = git.GetGitHandler().NewBranch(param.PathStore.OPPath, git.GetGitHandler().GenerateBranchName())
	if err != nil {
		return err
	}
	err = config.GetConfigHandler().ApplyMods()
	if err != nil {
		return err
	}
	err = git.GetGitHandler().AddBranch(param.PathStore.OPPath)
	if err != nil {
		return err
	}
	err = git.GetGitHandler().CommitBranch(param.PathStore.OPPath, config.GetConfigHandler().BuildModList())
	if err != nil {
		return err
	}
	return nil
}

func (installer *InstallerImpl) Reset() error {
	fmt.Println("Reseting...")
	err := git.GetGitHandler().CheckoutBranch(param.PathStore.OPPath, config.GetConfigHandler().GetConfig().OPVersion)
	if err != nil {
		return err
	}
	err = installer.RemoveAllOMMBranches()
	if err != nil {
		return err
	}
	err = file.GetFileHandler().RemoveFolder(param.PathStore.OMMPath)
	if err != nil {
		return err
	}
	_, err = config.GetConfigHandler().CreateConfig()
	if err != nil {
		return err
	}
	return nil
}

func (installer *InstallerImpl) RemoveAllOMMBranches() error {
	branches, err := git.GetGitHandler().ListBranch(param.PathStore.OPPath)
	if err != nil {
		return err
	}
	for _, b := range branches {
		if strings.Contains(b, "omm-") {
			err = git.GetGitHandler().RemoveBranch(param.PathStore.OPPath, b[2:])
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (installer *InstallerImpl) Remove(name string) error {
	fmt.Printf("Removing: %v", name)
	err := config.GetConfigHandler().RemoveMod(name)
	if err != nil {
		return err
	}
	err = file.GetFileHandler().RemoveFolder(filepath.Join(param.PathStore.OMMPath, name))
	if err != nil {
		return err
	}
	err = installer.Apply()
	if err != nil {
		return err
	}
	return nil
}

func (installer *InstallerImpl) Install(path string, force bool) error {
	if common.IsUrl(path) {
		err := installer.installFromUrl(path, force)
		if err != nil {
			return err
		}
	} else {
		err := installer.installFromFile(path, force)
		if err != nil {
			return err
		}
	}
	return installer.Apply()
}

func (installer *InstallerImpl) installFromFile(path string, force bool) error {
	man, err := manifest.GetManifestFromFile(path)
	if err != nil {
		return err
	}
	if mod, _ := config.GetConfigHandler().FindMod(man.Name); mod != nil && !force {
		fmt.Println("This Mod is already exist")
		return nil
	}
	fmt.Printf("Installing: %v@%v", man.Name, man.Version)
	err = config.GetConfigHandler().AddMod(&mod.Mod{
		Name:    man.Name,
		Version: man.Version,
		Url:     path,
	})
	if err != nil {
		return err
	}
	return nil
}
func (installer *InstallerImpl) installFromUrl(path string, force bool) error {
	specificVersion := ""
	if parts := strings.Split(path, "@"); len(parts) != 0 {
		specificVersion = parts[len(parts)-1]
		path = parts[0]
	}

	name, err := common.GetNameFromGithub(path)
	if err != nil {
		return err
	}

	if mod, _ := config.GetConfigHandler().FindMod(name); mod != nil {
		if !force {
			fmt.Println("This Mod is already exist")
			return nil
		} else {
			err = git.GetGitHandler().Pull(filepath.Join(param.PathStore.OMMPath, name))
			if err != nil {
				return err
			}
		}
	}
	err = cache.GetCacheHandler().Download(path)
	if err != nil {
		return err
	}
	version := ""
	if specificVersion == "" {
		version, err = getLatestModVersion(path)
		if err != nil {
			return err
		}
	} else {
		version = specificVersion
	}
	fmt.Printf("Installing: %v@%v", name, version)
	err = config.GetConfigHandler().AddMod(&mod.Mod{
		Name:    name,
		Version: version,
		Url:     path,
	})
	if err != nil {
		return err
	}
	return nil
}

func (installer *InstallerImpl) List() error {
	_, err := fmt.Println(config.GetConfigHandler().BuildModList())
	return err
}

func (installer *InstallerImpl) Init(OPPath string) error {
	_, err := file.GetFileHandler().LoadFile(filepath.Join(OPPath, config.CONFIG_FILE_NAME))
	if err != nil {
		e := errors.Unwrap(err)
		if _, ok := e.(*os.PathError); ok {
			fmt.Println("Initing...")
			version, err := git.GetGitHandler().GetBranchName(OPPath)
			if err != nil {
				return err
			}
			c := config.NewConfigHandler(version)
			return c.SaveConfig()
		} else {
			return err
		}
	}
	fmt.Println("You already have a config file.")
	return nil
}

func getLatestModVersion(rootPath string) (string, error) {
	branches, err := git.GetGitHandler().ListBranch(rootPath)
	if err != nil {
		return "", err
	}
	latestBranch := "0.0.0"
	for _, b := range branches {
		if strings.Contains(b, config.GetConfigHandler().GetConfig().OPVersion) {
			if semver.Compare(b, latestBranch) == 1 {
				latestBranch = b
			}
		}
	}
	if latestBranch == "0.0.0" {
		return "", errors.New("Cannot find compatable version of this mod")
	}
	return latestBranch, nil
}
