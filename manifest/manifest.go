package manifest

import "github.com/borgmon/openpilot-mod-manager/mod"

type Manifest struct {
	Name        string `yaml:"name"`
	DisplayName string `yaml:"displayName"`
	RepoUrl     string `yaml:"repoUrl"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Publisher   string `yaml:"publisher"`
}

func (man *Manifest) ToMod() *mod.Mod {
	return &mod.Mod{
		Name:    man.Name,
		Version: man.Version,
	}
}
