package cache

import "github.com/borgmon/openpilot-mod-manager/manifest"

type CacheHandler interface {
	GetManifest(name string) (*manifest.Manifest, error)
	GetManifests() ([]*manifest.Manifest, error)
	Download(url string) error
}
