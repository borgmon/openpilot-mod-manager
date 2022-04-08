package file

import (
	"github.com/borgmon/openpilot-mod-manager/patch"
)

type FileHandler interface {
	SaveFile(name string, data []byte) error
	LoadFile(name string) ([]byte, error)
	RemoveFile(name string) error
	AddLine(path string, m map[int]string) error
	ReplaceLine(path string, m map[int]string) error
	ListAllFilesRecursively(rootPath string) ([]string, error)
	ParsePatch(path string, opPath string) ([]patch.Patch, error)
	MoveFolderRecursively(move string, to string) error
	NewFileRecursively(filePath string) error
}
