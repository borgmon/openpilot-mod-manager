package file

type FileHandler interface {
	SaveFile(string, []byte) error
	LoadFile(string) ([]byte, error)
	RemoveFile(string) error
}
