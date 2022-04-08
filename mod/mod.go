package mod

import (
	"github.com/borgmon/openpilot-mod-manager/file"
)

type Mod struct {
	Name        string           `yaml:"name"`
	Version     string           `yaml:"version"`
	FileHandler file.FileHandler `yaml:"-"`
}
