package config

import "github.com/borgmon/openpilot-mod-manager/mod"

type ConfigLoader interface {
	CreateConfig() error
	RemoveConfig() error
	SaveConfig() error
	LoadConfig() (Config, error)
	AddMod(*mod.ModManifest) error
	RemoveMod(name string) error
	FindMod(name string) (*mod.ModManifest, error)
	SortMod() error
}
