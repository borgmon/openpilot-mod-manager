package config

import (
	"path/filepath"

	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/injector"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/borgmon/openpilot-mod-manager/mod"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ConfigHandlerImpl struct {
	Config     *Config
	ConfigPath string
	CachePath  string
	OPPath     string
}

func NewConfigHandler(ConfigPath string, CachePath string, OPPath string) ConfigHandler {
	return &ConfigHandlerImpl{
		Config:     &Config{OPVersion: "master"},
		ConfigPath: ConfigPath,
		CachePath:  CachePath,
		OPPath:     OPPath,
	}
}

func (config ConfigHandlerImpl) CreateConfig() (*Config, error) {
	config.Config = &Config{Mods: []*mod.Mod{}}
	err := config.SaveConfig()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return config.Config, nil
}

func (config ConfigHandlerImpl) RemoveConfig() error {
	return file.GetFileHandler().RemoveFile(config.ConfigPath)
}

func (config ConfigHandlerImpl) SaveConfig() error {
	bytes, err := yaml.Marshal(config.Config)
	if err != nil {
		return errors.WithStack(err)
	}
	err = file.GetFileHandler().SaveFile(config.ConfigPath, bytes)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (config ConfigHandlerImpl) LoadConfig() (*Config, error) {
	bytes, err := file.GetFileHandler().LoadFile(config.ConfigPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = yaml.Unmarshal(bytes, config.Config)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return config.Config, nil
}

func (config ConfigHandlerImpl) AddMod(mod *mod.Mod) error {
	config.Config.Mods = append(config.Config.Mods, mod)
	err := config.SaveConfig()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (config ConfigHandlerImpl) RemoveMod(name string) error {
	mods := []*mod.Mod{}

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

func (config ConfigHandlerImpl) FindMod(name string) (*mod.Mod, error) {
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

func (config ConfigHandlerImpl) GetManifests() ([]*manifest.Manifest, error) {
	result := []*manifest.Manifest{}
	for _, mod := range config.Config.Mods {
		man, err := config.GetManifest(mod.Name)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result = append(result, man)
	}
	return result, nil
}

func (config ConfigHandlerImpl) ApplyMods() error {
	for _, mod := range config.Config.Mods {
		rootPath := filepath.Join(config.CachePath, mod.Name)
		paths, err := file.GetFileHandler().ListAllFilesRecursively(rootPath)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, path := range paths {
			absPath := filepath.Join(config.OPPath, path)
			patches, err := file.GetFileHandler().ParsePatch(path, absPath)
			if err != nil {
				return errors.WithStack(err)
			}
			for _, p := range patches {
				injector.GetInjector().Pending(p)
			}
		}
	}
	injector.GetInjector().Inject()
	return nil
}

func (config ConfigHandlerImpl) GetManifest(name string) (*manifest.Manifest, error) {
	mod, err := config.FindMod(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	path := filepath.Join(config.CachePath, mod.Name, manifest.MANIFEST_FILE_NAME)
	data, err := file.GetFileHandler().LoadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	man := &manifest.Manifest{}
	err = yaml.Unmarshal(data, man)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return man, nil

}

func (config ConfigHandlerImpl) GetConfig() *Config {
	return config.Config
}
