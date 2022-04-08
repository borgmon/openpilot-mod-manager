package source

import "github.com/borgmon/openpilot-mod-manager/manifest"

type Source interface {
	DownloadToCache() (*manifest.Manifest, error)
	GetName() (string, error)
}
