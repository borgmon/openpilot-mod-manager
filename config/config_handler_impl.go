package config

import (
	"path/filepath"
	"strings"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/git"
	"github.com/borgmon/openpilot-mod-manager/injector"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/borgmon/openpilot-mod-manager/mod"
	"github.com/borgmon/openpilot-mod-manager/param"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ConfigHandlerImpl struct {
	Config *Config
}

var configHandlerInstance ConfigHandler

func NewConfigHandler(OPVersion string) ConfigHandler {
	c := &ConfigHandlerImpl{Config: &Config{
		OPVersion: OPVersion,
		Mods:      []*mod.Mod{},
	}}
	configHandlerInstance = c
	return c
}

func LoadConfigHandler() (ConfigHandler, error) {
	c := &ConfigHandlerImpl{Config: &Config{
		OPVersion: "",
		Mods:      []*mod.Mod{},
	}}
	_, err := c.LoadConfig()
	if err != nil {
		return nil, err
	}
	configHandlerInstance = c
	return c, nil
}

func GetConfigHandler() ConfigHandler {
	return configHandlerInstance
}

func (config *ConfigHandlerImpl) CreateConfig() (*Config, error) {
	c := NewConfigHandler(config.Config.OPVersion)
	err := c.SaveConfig()
	if err != nil {
		return nil, err
	}
	return c.GetConfig(), nil
}

func (config *ConfigHandlerImpl) RemoveConfig() error {
	return file.GetFileHandler().RemoveFile(param.PathStore.ConfigPath)
}

func (config *ConfigHandlerImpl) SaveConfig() error {
	bytes, err := yaml.Marshal(config.Config)
	if err != nil {
		return errors.WithStack(err)
	}
	err = file.GetFileHandler().SaveFile(param.PathStore.ConfigPath, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (config *ConfigHandlerImpl) LoadConfig() (*Config, error) {
	bytes, err := file.GetFileHandler().LoadFile(param.PathStore.ConfigPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, config.Config)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return config.Config, nil
}

func (config *ConfigHandlerImpl) AddMod(mod *mod.Mod) error {
	if r, _ := config.FindMod(mod.Name); r != nil {
		return nil
	}
	config.Config.Mods = append(config.Config.Mods, mod)
	err := config.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func (config *ConfigHandlerImpl) RemoveMod(name string) error {
	mods := []*mod.Mod{}

	for _, mod := range config.Config.Mods {
		if mod.Name != name {
			mods = append(mods, mod)
		}
	}

	config.Config.Mods = mods
	err := config.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func (config *ConfigHandlerImpl) FindMod(name string) (*mod.Mod, error) {
	for _, mod := range config.Config.Mods {
		if mod.Name == name {
			return mod, nil
		}
	}
	return nil, nil
}

// TODO: dependencies
func (config ConfigHandlerImpl) SortMod() error {
	return errors.New("SortMod not implemented")
}

func (config *ConfigHandlerImpl) ApplyMods() error {
	for _, mod := range config.Config.Mods {
		rootPath := ""
		if common.IsUrl(mod.Url) {
			rootPath = filepath.Join(param.PathStore.OMMPath, mod.Name)
			err := git.GetGitHandler().CheckoutBranch(rootPath, mod.Version)
			if err != nil {
				return err
			}

		} else {
			rootPath = mod.Url
		}

		paths, err := file.GetFileHandler().ListAllFilesRecursively(rootPath)
		if err != nil {
			return err
		}
		paths = filterFiles(paths)
		for _, path := range paths {
			relativePath := strings.ReplaceAll(path, rootPath, "")
			absPath := filepath.Join(param.PathStore.OPPath, relativePath)
			patches, err := file.GetFileHandler().ParsePatch(path, absPath)
			if err != nil {
				return err
			}
			for _, p := range patches {
				p.Mod = mod
				injector.GetInjector().Pending(p)
			}
		}
	}
	injector.GetInjector().Inject()
	return nil
}

func filterFiles(paths []string) []string {
	blackList := []string{manifest.MANIFEST_FILE_NAME, CONFIG_FILE_NAME, ".git"}
	results := []string{}
	for _, path := range paths {
		if !pathInBlackList(path, blackList) {
			results = append(results, path)
		}
	}
	return results
}

func pathInBlackList(path string, blackList []string) bool {
	for _, part := range strings.Split(path, "/") {
		if partsInBlackList(part, blackList) {
			return true
		}
	}
	return false
}

func partsInBlackList(path string, blackList []string) bool {
	for _, black := range blackList {
		if black == path {
			return true
		}
	}
	return false
}

func (config *ConfigHandlerImpl) GetConfig() *Config {
	return config.Config
}

func (config *ConfigHandlerImpl) BuildModList() string {
	result := []string{}
	for _, m := range config.Config.Mods {
		result = append(result, m.Name+":"+m.Version)
	}
	return strings.Join(result, "\n")
}
