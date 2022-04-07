package file

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type FileHandlerImpl struct{}

func (handler *FileHandlerImpl) SaveFile(name string, data []byte) error {
	err := ioutil.WriteFile("temp.txt", data, 0666)
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
