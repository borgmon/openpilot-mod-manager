package injector

import "github.com/borgmon/openpilot-mod-manager/mod"

type Injector interface {
	Pending([]mod.ModManifest) error
	inject() error
}
