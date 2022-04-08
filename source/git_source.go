package source

import (
	"path/filepath"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/git"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type GitSource struct {
	RemoteUrl  string
	GitHandler git.GitHandler
	CachePath  string
}

func (source *GitSource) DownloadToCache() (*manifest.Manifest, error) {
	err := source.GitHandler.Clone(source.RemoteUrl)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	name, err := source.GetName()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	man, err := source.getManifest(filepath.Join(source.CachePath, name))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return man, nil
}

func (source *GitSource) getManifest(localPath string) (*manifest.Manifest, error) {
	data, err := file.GetFileHandler().LoadFile(localPath)
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
func (source *GitSource) GetName() (string, error) {
	name, err := common.GetProjectFromGithub(source.RemoteUrl)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return name, nil
}
