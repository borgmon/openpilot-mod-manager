package manifest

import (
	"path/filepath"

	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ManifestHandlerImpl struct{}

func GetManifestHandler() ManifestHandler {
	return &ManifestHandlerImpl{}
}

func (handler *ManifestHandlerImpl) Init(path string) error {
	manifest := &Manifest{
		Name:        "my-mod",
		DisplayName: "my mod",
		RepoUrl:     "https://github.com/myname/my-mod",
		Version:     "v0.8.13-1",
		Description: "This is my mod",
		Publisher:   "my name",
	}
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return errors.WithStack(err)
	}
	return file.GetFileHandler().SaveFile(filepath.Join(path, MANIFEST_FILE_NAME), data)
}
