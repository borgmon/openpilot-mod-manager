package patch

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/borgmon/openpilot-mod-manager/mod"
	"github.com/pkg/errors"
)

type Patch struct {
	Path       string
	LineNumber int
	Data       string
	Operand    string
	Mod        *mod.Mod
}

const (
	TypeOperandReplace = "<<<"
	TypeOperandAppend  = ">>>"
)

// key format: relative_path/filename.ext#line_num
func KeyToPatch(key string, data string) (*Patch, error) {
	paths := strings.Split(key, "#")
	lineNum, err := strconv.Atoi(paths[1])
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Patch{
		Path:       paths[0],
		LineNumber: lineNum,
		Data:       data,
	}, nil
}

func (patch *Patch) ToKey() string {
	return fmt.Sprintf("%v#%v\n", patch.Path, patch.LineNumber)
}

func (patch *Patch) AppendData(data string) {
	patch.Data += "\n" + data
}
