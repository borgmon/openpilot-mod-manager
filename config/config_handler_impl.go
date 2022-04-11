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
	"github.com/borgmon/openpilot-mod-manager/version"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gopkg.in/yaml.v2"
)

var filePathBlackList = []string{manifest.MANIFEST_FILE_NAME, CONFIG_FILE_NAME, ".git"}

type ConfigHandlerImpl struct {
	Config *Config
}

var configHandlerInstance ConfigHandler

func NewConfigHandler(OPVersion string) ConfigHandler {
	c := &ConfigHandlerImpl{Config: &Config{
		OPVersion:  OPVersion,
		OMMVersion: version.OMMVersion,
		Mods:       []*mod.Mod{},
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
	if _, ok := config.FindMod(mod.Name); ok {
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
	config.Config.Mods = lo.Filter(config.Config.Mods, func(m *mod.Mod, i int) bool {
		return m.Name != name
	})
	err := config.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func (config *ConfigHandlerImpl) FindMod(name string) (*mod.Mod, bool) {
	return lo.Find(config.Config.Mods, func(m *mod.Mod) bool {
		return m.Name == name
	})
}

// TODO: dependencies
func (config ConfigHandlerImpl) SortMod() error {
	if baseMod, ok := config.FindMod(param.BaseModName); ok {
		config.Config.Mods = append([]*mod.Mod{baseMod}, lo.Filter(config.Config.Mods, func(mod *mod.Mod, i int) bool {
			return mod.Name != param.BaseModName
		})...)
	} else {
		return errors.New("No base mod found")
	}

	return nil
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
		paths = filterBlacklist(paths, filePathBlackList)
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

func filterBlacklist(paths, blackList []string) []string {
	return lo.Filter(paths, func(path string, i int) bool {
		parts := strings.Split(path, "/")
		_, ok := lo.Find(parts, func(t string) bool {
			return lo.Contains(blackList, t)
		})
		return !ok
	})
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
