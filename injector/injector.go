package injector

import "github.com/borgmon/openpilot-mod-manager/patch"

type Injector interface {
	Pending(p patch.Patch) error
	Inject()
}
