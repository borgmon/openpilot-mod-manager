package config

import (
	"github.com/borgmon/openpilot-mod-manager/mod"
)

type Config struct {
	Mods []*mod.ModManifest `yaml:"mods"`
}
