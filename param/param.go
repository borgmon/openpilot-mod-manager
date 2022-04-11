package param

type Path struct {
	ConfigPath string
	OMMPath    string
	OPPath     string
}

type Config struct {
	Verbose bool
}

var PathStore *Path
var ConfigStore *Config

const BaseModUrl = "https://github.com/borgmon/omm-base@main"
const BaseModName = "omm-base"

func NewParam(ConfigPath string, OMMPath string, OPPath string, Verbose bool) {
	PathStore = &Path{
		ConfigPath: ConfigPath,
		OMMPath:    OMMPath,
		OPPath:     OPPath,
	}
	ConfigStore = &Config{
		Verbose: Verbose,
	}
}
