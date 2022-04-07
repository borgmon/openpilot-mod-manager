package config

import (
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/mod"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ConfigHandlerImpl struct {
	Config      *Config
	FileHandler file.FileHandler
	configName  string
}

func (config ConfigHandlerImpl) CreateConfig() error {
	config.Config.Mods = []*mod.ModManifest{}
	err := config.SaveConfig()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (config ConfigHandlerImpl) RemoveConfig() error {
	return config.FileHandler.RemoveFile(config.configName)
}

func (config ConfigHandlerImpl) SaveConfig() error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return errors.WithStack(err)
	}
	err = config.FileHandler.SaveFile(config.configName, bytes)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (config ConfigHandlerImpl) LoadConfig() (*Config, error) {
	bytes, err := config.FileHandler.LoadFile(config.configName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = yaml.Unmarshal(bytes, config.Config)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return config.Config, nil
}

func (config ConfigHandlerImpl) AddMod(mod *mod.ModManifest) error {
	config.Config.Mods = append(config.Config.Mods, mod)
	err := config.SaveConfig()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (config ConfigHandlerImpl) RemoveMod(name string) error {
	mods := []*mod.ModManifest{}
	for _, mod := range config.Config.Mods {
		if mod.Name != name {
			mods = append(mods, mod)
		}
	}

	config.Config.Mods = mods
	err := config.SaveConfig()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (config ConfigHandlerImpl) FindMod(name string) *mod.ModManifest {
	for _, mod := range config.Config.Mods {
		if mod.Name == name {
			return mod
		}
	}
	return nil
}

// TODO: dependencies
func (config ConfigHandlerImpl) SortMod() error {
	return nil
}
