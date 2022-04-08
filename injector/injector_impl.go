package injector

import (
	"strings"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/file"
	ommerrors "github.com/borgmon/openpilot-mod-manager/omm-errors"
	"github.com/borgmon/openpilot-mod-manager/patch"
	"github.com/pkg/errors"
)

type InjectorImpl struct {
	Changes map[string]patch.Patch // [filepath#linenum]:patch
}

var injectorInstance Injector

func GetInjector() Injector {
	if injectorInstance != nil {
		return injectorInstance
	}
	return &InjectorImpl{Changes: map[string]patch.Patch{}}
}

func (injector *InjectorImpl) Pending(p patch.Patch) error {
	if _, ok := injector.Changes[p.ToKey()]; ok {
		if p.GetOperand() == patch.TypeOperandReplace {
			return ommerrors.NewReplaceConflictError(p.GetFilePath(), p.GetLineNumber())
		}

		injector.Changes[p.ToKey()].AppendData(p.GetData())
	} else {
		injector.Changes[p.ToKey()] = p
	}
	return nil
}

func (injector *InjectorImpl) Inject() {
	// remap into [filepath]:[]patch
	appendMap := map[string][]patch.Patch{}
	replaceMap := map[string][]patch.Patch{}
	for k := range injector.Changes {
		parts := strings.Split(k, "#")
		path := parts[0]

		switch injector.Changes[k].GetOperand() {
		case patch.TypeOperandAppend:
			if v, ok := appendMap[path]; ok {
				appendMap[path] = append(v, injector.Changes[k])
			} else {
				appendMap[path] = []patch.Patch{injector.Changes[k]}
			}
		case patch.TypeOperandReplace:
			if v, ok := replaceMap[path]; ok {
				replaceMap[path] = append(v, injector.Changes[k])
			} else {
				replaceMap[path] = []patch.Patch{injector.Changes[k]}
			}
		}

	}
	for k := range appendMap {
		injector.doInject(k, appendMap[k], replaceMap[k])
	}
}

func (injector *InjectorImpl) doInject(path string, appends []patch.Patch, replaces []patch.Patch) error {
	// remap into [line num]:patch
	appendMap := map[int]string{}
	// replaceMap := map[int]string{}
	for _, patch := range appends {
		appendMap[patch.GetLineNumber()] = patch.GetData()
	}
	// for _, patch := range replaces {
	// 	replaceMap[patch.GetLineNumber()] = patch.GetData()
	// }
	err := file.GetFileHandler().AddLine(path, appendMap)
	if err != nil {
		return common.LogIfErr(errors.WithStack(err))
	}
	// err = file.GetFileHandler().ReplaceLine(path, replaceMap)
	// if err != nil {
	// 	return common.LogIfErr(errors.WithStack(err))
	// }
	return nil
}
