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
