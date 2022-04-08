package source

import (
	"path/filepath"

	"github.com/borgmon/openpilot-mod-manager/config"
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/borgmon/openpilot-mod-manager/param"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type LocalSource struct {
	LocalPath     string
	ConfigHandler config.ConfigHandler
}

func (source *LocalSource) DownloadToCache() (*manifest.Manifest, error) {
	man, err := source.getManifest()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = file.GetFileHandler().CopyFolderRecursively(source.LocalPath, param.PathStore.OMMPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return man, nil
}

func (source *LocalSource) getManifest() (*manifest.Manifest, error) {
	data, err := file.GetFileHandler().LoadFile(filepath.Join(source.LocalPath, manifest.MANIFEST_FILE_NAME))
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
