package file

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/borgmon/openpilot-mod-manager/patch"
	"github.com/pkg/errors"
)

type FileHandlerImpl struct{}

func GetFileHandler() FileHandler {
	return &FileHandlerImpl{}
}

func (handler *FileHandlerImpl) SaveFile(name string, data []byte) error {
	err := ioutil.WriteFile(name, data, 0666)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *FileHandlerImpl) LoadFile(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

func (handler *FileHandlerImpl) RemoveFile(name string) error {
	return os.Remove(name)
}

func (handler *FileHandlerImpl) AddLine(path string, m map[int]string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		err := handler.NewFile(path)
		if err != nil {
			return errors.WithStack(err)
		}
		data, err = ioutil.ReadFile(path)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	text := string(data)
	lines := strings.Split(text, "\n")
	newLines := lines
	offset := 0
	for i, _ := range lines {
		if appendText, ok := m[i]; ok {
			l := []string{appendText}
			newLines = append(newLines[:i+offset], append(l, newLines[(i+offset):]...)...)
			offset++
		}
	}

	err = handler.SaveFile(path, []byte(strings.Join(newLines, "\n")))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *FileHandlerImpl) ReplaceLine(path string, m map[int]string) error {
	return errors.New("<<< not implement yet")
}

func (handler *FileHandlerImpl) ListAllFilesRecursively(rootPath string) ([]string, error) {
	result := []string{}
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.WithStack(err)
			}
			if !info.IsDir() {
				result = append(result, path)
			}
			return nil
		})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (handler *FileHandlerImpl) CopyFolderRecursively(move string, to string) error {
	err := execRunner("cp", "-R", move, to)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *FileHandlerImpl) NewFolder(path string) error {
	err := execRunner("mkdir", "-p", path)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *FileHandlerImpl) NewBranch(path string) error {
	err := execRunner("mkdir", "-p", path)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *FileHandlerImpl) NewFile(path string) error {
	touch := exec.Command("touch", path)
	err := touch.Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func execRunner(name string, arg ...string) error {
	e := exec.Command(name, arg...)
	b := new(strings.Builder)
	e.Stdout = b
	err := e.Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *FileHandlerImpl) NewFileRecursively(filePath string) error {
	path := common.GetPathFromFilePath(filePath)

	err := handler.NewFolder(path)
	if err != nil {
		return errors.WithStack(err)
	}
	err = handler.NewFile(filePath)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *FileHandlerImpl) RemoveFolder(path string) error {
	return os.RemoveAll(path)
}

func (handler *FileHandlerImpl) ParsePatch(path string, opPath string) ([]patch.Patch, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	i := 1
	result := []patch.Patch{}
	operand := ""
	buf := ""
	start := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				result = append(result, &patch.PatchImpl{
					Path:       opPath,
					LineNumber: start,
					Data:       buf,
					Operand:    operand,
				})
				return result, nil
			}
			return nil, errors.WithStack(err)
		}
		if op := getOperands(line); op != "" {
			if buf != "" {
				result = append(result, &patch.PatchImpl{
					Path:       opPath,
					LineNumber: start,
					Data:       buf,
					Operand:    operand,
				})
			}
			operand = op
			buf = ""
			start, err = parseLineNum(line)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		} else {
			buf += line
		}
		i++
	}

}

func parseLineNum(line string) (int, error) {
	parts := strings.Split(line, "#")
	last := parts[len(parts)-1]
	last = last[:len(last)-1]
	num, err := strconv.Atoi(last)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return num, nil
}

func getOperands(line string) string {
	if strings.Contains(line, patch.TypeOperandAppend) {
		return patch.TypeOperandAppend
	}
	if strings.Contains(line, patch.TypeOperandReplace) {
		return patch.TypeOperandReplace
	}
	return ""
}
