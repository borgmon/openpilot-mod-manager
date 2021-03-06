package file

import (
	"github.com/borgmon/openpilot-mod-manager/patch"
)

type FileHandler interface {
	SaveFile(name string, data []byte) error
	LoadFile(name string) ([]byte, error)
	RemoveFile(name string) error
	ModifyFile(path string, patchMap map[int]string, deleteMap map[int]string) error
	ListAllFilesRecursively(rootPath string) ([]string, error)
	ParsePatch(path string, opPath string) ([]*patch.Patch, error)
	CopyFolderRecursively(move string, to string) error
	NewFileRecursively(filePath string) error
	NewFolder(path string) error
	NewFile(path string) error
	RemoveFolder(path string) error
}
