package cache

import (
	"path/filepath"

	"github.com/borgmon/openpilot-mod-manager/config"
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
	mod, err := config.GetConfigHandler().FindMod(name)
	if err != nil {
		return nil, err
	}
	man, err := manifest.GetManifestFromFile(filepath.Join(param.PathStore.OMMPath, mod.Name))
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

func (cache *CacheHandlerImpl) Download(url string) error {
	err := git.GetGitHandler().Clone(param.PathStore.OMMPath, url)
	if err != nil {
		return err
	}
	return nil
}
