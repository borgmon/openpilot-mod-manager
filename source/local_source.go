package source

import (
	"path/filepath"

	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type LocalSource struct {
	LocalPath string
	CachePath string
}

func (source *LocalSource) DownloadToCache() (*manifest.Manifest, error) {
	man, err := source.getManifest()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = file.GetFileHandler().MoveFolderRecursively(source.LocalPath, filepath.Join(source.CachePath, man.Name))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return man, nil
}

func (source *LocalSource) getManifest() (*manifest.Manifest, error) {
	data, err := file.GetFileHandler().LoadFile(source.LocalPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	man := &manifest.Manifest{}
	err = yaml.Unmarshal(data, man)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return man, nil
}

func (source *LocalSource) GetName() (string, error) {
	man, err := source.getManifest()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return man.Name, nil
}
