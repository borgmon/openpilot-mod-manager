package manifest

import (
	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Name        string `yaml:"name"`
	DisplayName string `yaml:"displayName"`
	RepoUrl     string `yaml:"repoUrl"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Publisher   string `yaml:"publisher"`
}

const MANIFEST_FILE_NAME = "manifest.yml"

func (man *Manifest) BuildInfo() (string, error) {
	out, err := yaml.Marshal(man)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
