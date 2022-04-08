package manifest

type ManifestHandler interface {
	Init(path string) error
}
