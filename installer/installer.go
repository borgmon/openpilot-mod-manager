package installer

type Installer interface {
	Reset() error
	Apply() error
	Install(path string, force bool) error
	Remove(name string) error
	List() error
	Init(OPPath string) error
	Info(name string) error
}
