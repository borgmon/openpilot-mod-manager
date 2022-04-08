package source

import (
	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/config"
	"github.com/borgmon/openpilot-mod-manager/git"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/pkg/errors"
)

type GitSource struct {
	RemoteUrl     string
	ConfigHandler config.ConfigHandler
}

func (source *GitSource) DownloadToCache() (*manifest.Manifest, error) {
	err := git.GetGitHandler().Clone(source.ConfigHandler.GetPaths().CachePath, source.RemoteUrl)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	name, err := source.GetName()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	man, err := source.ConfigHandler.GetManifest(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return man, nil
}

func (source *GitSource) GetName() (string, error) {
	name, err := common.GetNameFromGithub(source.RemoteUrl)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return name, nil
}
