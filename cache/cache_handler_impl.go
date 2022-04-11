package cache

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/config"
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/git"
	"github.com/borgmon/openpilot-mod-manager/manifest"
	"github.com/borgmon/openpilot-mod-manager/param"
)

type CacheHandlerImpl struct{}

var cacheHandlerInstance CacheHandler

func GetCacheHandler() CacheHandler {
	if cacheHandlerInstance != nil {
		return cacheHandlerInstance
	}
	cacheHandlerInstance = &CacheHandlerImpl{}
	return cacheHandlerInstance
}

func (cache *CacheHandlerImpl) GetManifest(name string) (*manifest.Manifest, error) {
	man, err := manifest.GetManifestFromFile(filepath.Join(param.PathStore.OMMPath, name))
	if err != nil {
		return nil, err
	}
	return man, nil
}

func (cache *CacheHandlerImpl) GetManifests() ([]*manifest.Manifest, error) {
	result := []*manifest.Manifest{}
	for _, mod := range config.GetConfigHandler().GetConfig().Mods {
		man, err := cache.GetManifest(mod.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, man)
	}
	return result, nil
}

func (cache *CacheHandlerImpl) Download(url string, force bool) error {
	name, err := common.GetNameFromGithub(url)
	if err != nil {
		return err
	}
	path := filepath.Join(param.PathStore.OMMPath, name)
	_, err = cache.GetManifest(name)
	if err != nil {
		e := errors.Unwrap(err)
		if _, ok := e.(*os.PathError); ok {
			err := git.GetGitHandler().Clone(param.PathStore.OMMPath, url)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if force {
		err = file.GetFileHandler().RemoveFolder(path)
		if err != nil {
			return err
		}
		err = git.GetGitHandler().Clone(param.PathStore.OMMPath, url)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cache *CacheHandlerImpl) RemoveCache() error {
	return file.GetFileHandler().RemoveFolder(param.PathStore.OMMPath)
}
