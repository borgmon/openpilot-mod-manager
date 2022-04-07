package mod

import (
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ModHandlerImpl struct {
	FileHandler  file.FileHandler
	TemplatePath string
	CurrentPath  string
}

func (handler *ModHandlerImpl) Init() error {
	manifest := &ModManifest{
		Name:        "my-mod",
		DisplayName: "my mod",
		RepoUrl:     "https://github.com/myname/my-mod",
		Version:     .01,
		Publisher:   "my name",
	}
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return errors.WithStack(err)
	}
	return handler.FileHandler.SaveFile(handler.CurrentPath, data)
}
