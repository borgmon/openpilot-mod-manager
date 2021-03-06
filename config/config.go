package config

import "github.com/borgmon/openpilot-mod-manager/mod"

type Config struct {
	OPVersion  string     `yaml:"OPVersion"`
	OMMVersion string     `yaml:"OMMVersion"`
	Mods       []*mod.Mod `yaml:"mods"`
}

const CONFIG_FILE_NAME = "omm.yml"
const CACHEPATH = ".omm"
