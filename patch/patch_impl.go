package patch

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type PatchImpl struct {
	Path       string
	LineNumber int
	Data       string
	Operand    string
}

// key format: relative_path/filename.ext#line_num
func KeyToPatch(key string, data string) (*PatchImpl, error) {
	paths := strings.Split(key, "#")
	lineNum, err := strconv.Atoi(paths[1])
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &PatchImpl{
		Path:       paths[0],
		LineNumber: lineNum,
		Data:       data,
	}, nil
}

func (patch *PatchImpl) ToKey() string {
	return patch.Path + "#" + strconv.Itoa(patch.LineNumber)
}

func (patch *PatchImpl) GetLineNumber() int {
	return patch.LineNumber
}

func (patch *PatchImpl) GetData() string {
	return patch.Data
}

func (patch *PatchImpl) GetFilePath() string {
	return patch.Path
}

func (patch *PatchImpl) GetOperand() string {
	return patch.Operand
}

func (patch *PatchImpl) AppendData(data string) {
	patch.Data += "\n" + data
}
