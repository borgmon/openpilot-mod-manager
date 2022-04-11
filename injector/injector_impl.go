package injector

import (
	"fmt"
	"strings"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/file"
	"github.com/borgmon/openpilot-mod-manager/patch"
)

type InjectorImpl struct {
	Changes map[string]*patch.Patch // [filepath#linenum]:patch
}

var injectorInstance Injector

func GetInjector() Injector {
	if injectorInstance != nil {
		return injectorInstance
	}
	injectorInstance = &InjectorImpl{
		Changes: map[string]*patch.Patch{},
	}
	return injectorInstance
}

func (injector *InjectorImpl) Pending(p *patch.Patch) error {
	fmt.Printf("Pending patch: mod=%v\tfile=%v\tmode=%v\n", p.Mod.Name, p.ToKey(), p.Operand)
	if _, ok := injector.Changes[p.ToKey()]; ok {
		switch injector.Changes[p.ToKey()].Operand {
		case patch.TypeOperandDelete:
			if p.Operand == patch.TypeOperandDelete {
				return nil
			} else {
				injector.Changes[p.ToKey()].AppendData(p.Data)
			}
		case patch.TypeOperandAdd:
			if p.Operand == patch.TypeOperandDelete {
				injector.Changes[p.ToKey()].Operand = patch.TypeOperandDelete
			} else {
				injector.Changes[p.ToKey()].AppendData(p.Data)
			}
		}
	} else {
		injector.Changes[p.ToKey()] = p
	}
	return nil
}

func (injector *InjectorImpl) Inject() {
	// remap into [filepath]:[]patch
	patchMap := map[string][]*patch.Patch{}
	for k := range injector.Changes {
		parts := strings.Split(k, "#")
		path := parts[0]

		if v, ok := patchMap[path]; ok {
			patchMap[path] = append(v, injector.Changes[k])
		} else {
			patchMap[path] = []*patch.Patch{injector.Changes[k]}
		}

	}
	for k := range patchMap {
		injector.doInject(k, patchMap[k])
	}
}

func (injector *InjectorImpl) doInject(path string, patchMap []*patch.Patch) error {
	// remap into [line num]:patch
	appendMap := map[int]string{}
	deleteMap := map[int]string{}
	for _, p := range patchMap {
		fmt.Printf("Inject patch: mod=%v\tfile=%v\tmode=%v\n", p.Mod.Name, p.ToKey(), p.Operand)
		switch p.Operand {
		case patch.TypeOperandAdd:
			appendMap[p.LineNumber] = p.Data
		case patch.TypeOperandDelete:
			deleteMap[p.LineNumber] = p.Data
		}
	}
	err := file.GetFileHandler().ModifyFile(path, appendMap, deleteMap)
	if err != nil {
		return common.LogIfErr(err)
	}
	return nil
}
