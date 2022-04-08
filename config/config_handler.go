package config

import (
	"github.com/borgmon/openpilot-mod-manager/mod"
)

type ConfigHandler interface {
	CreateConfig() error
	RemoveConfig() error
	SaveConfig() error
	LoadConfig() (*Config, error)
	AddMod(mod *mod.Mod) error
	RemoveMod(name string) error
	FindMod(name string) (*mod.Mod, error)
	SortMod() error
	ApplyMods() error
	// ListManifests() ([]*manifest.Manifest, error)
}
