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

var fileHandlerInstance FileHandler

func GetFileHandler() FileHandler {
	if fileHandlerInstance != nil {
		return fileHandlerInstance
	}
	fileHandlerInstance = &FileHandlerImpl{}
	return fileHandlerInstance
}

func (handler *FileHandlerImpl) SaveFile(name string, data []byte) error {
	return errors.WithStack(ioutil.WriteFile(name, data, 0666))
}

func (handler *FileHandlerImpl) LoadFile(name string) ([]byte, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return b, nil
}

func (handler *FileHandlerImpl) RemoveFile(name string) error {
	return errors.WithStack(os.Remove(name))
}

func (handler *FileHandlerImpl) ModifyFile(path string, addMap map[int]string, deleteMap map[int]string) error {
	data, err := handler.LoadFile(path)
	if err != nil {
		e := errors.Unwrap(err)
		if _, ok := e.(*os.PathError); ok {
			err := handler.NewFile(path)
			if err != nil {
				return err
			}
			data, err = handler.LoadFile(path)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	text := string(data)
	lines := strings.Split(text, "\n")
	newLines := lines
	offset := 0
	for i, _ := range lines {
		if appendText, ok := addMap[i]; ok {
			l := []string{appendText}
			newLines = append(newLines[:i+offset-1], append(l, newLines[(i+offset-1):]...)...)
			offset++
		}
		if replaceText, ok := deleteMap[i]; ok {
			l := []string{replaceText}
			if replaceText != "" {
				newLines = append(newLines[:i+offset-1], append(l, newLines[(i+offset):]...)...)
			} else {
				newLines = append(newLines[:i+offset-1], newLines[(i+offset):]...)
				offset--
			}

		}
	}
	if appendText, ok := addMap[len(lines)]; ok {
		l := []string{appendText}
		newLines = append(newLines, l...)
	}
	if replaceText, ok := deleteMap[len(lines)]; ok {
		l := []string{replaceText}
		if replaceText != "" {
			newLines = append(newLines[:len(newLines)-1], l...)
		} else {
			newLines = newLines[:len(newLines)-1]
		}
	}

	return handler.SaveFile(path, []byte(strings.Join(newLines, "\n")))
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
	return errors.WithStack(execRunner("cp", "-R", move, to))
}

func (handler *FileHandlerImpl) NewFolder(path string) error {
	return errors.WithStack(execRunner("mkdir", "-p", path))
}

func (handler *FileHandlerImpl) NewFile(path string) error {
	return errors.WithStack(execRunner("touch", path))
}

func execRunner(name string, arg ...string) error {
	e := exec.Command(name, arg...)
	b := new(strings.Builder)
	e.Stdout = b
	common.LogIfVarbose(b.String())
	return errors.WithStack(e.Run())
}

func (handler *FileHandlerImpl) NewFileRecursively(filePath string) error {
	path := common.GetPathFromFilePath(filePath)

	err := handler.NewFolder(path)
	if err != nil {
		return err
	}
	return handler.NewFile(filePath)
}

func (handler *FileHandlerImpl) RemoveFolder(path string) error {
	return errors.WithStack(os.RemoveAll(path))
}

func (handler *FileHandlerImpl) ParsePatch(path string, opPath string) ([]*patch.Patch, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	i := 1
	result := []*patch.Patch{}
	operand := ""
	buf := ""
	start := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if operand != "" {
					result = append(result, &patch.Patch{
						Path:       opPath,
						LineNumber: start,
						Data:       removeTrailingNewLine(buf),
						Operand:    operand,
					})
					return result, nil
				} else {
					return result, nil
				}
			}
			return nil, errors.WithStack(err)
		}
		if op := getOperands(line); op != "" {
			if buf != "" || operand == patch.TypeOperandDelete {
				result = append(result, &patch.Patch{
					Path:       opPath,
					LineNumber: start,
					Data:       removeTrailingNewLine(buf),
					Operand:    operand,
				})
			}
			operand = op
			buf = ""
			start, err = parseLineNum(line)
			if err != nil {
				return nil, err
			}
		} else {
			buf += line
		}
		i++
	}

}

func removeTrailingNewLine(str string) string {
	if str != "" && str[len(str)-1] == '\n' {
		return str[:len(str)-1]
	} else {
		return str
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
	if strings.Contains(line, patch.TypeOperandAdd) {
		return patch.TypeOperandAdd
	}
	if strings.Contains(line, patch.TypeOperandDelete) {
		return patch.TypeOperandDelete
	}
	return ""
}
