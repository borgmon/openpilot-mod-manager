package common

import (
	"fmt"

	"github.com/borgmon/openpilot-mod-manager/param"
)

func PanicIfErr(err error) error {
	if err != nil {
		panic(err)
	}
	return err
}

func LogIfErr(err error) error {
	if err != nil && param.ConfigStore.Verbose {
		fmt.Printf("%+v", err)
	}
	return err
}

func LogIfVarbose(str string) {
	if param.ConfigStore.Verbose {
		fmt.Printf("%+v", str)
	}
}
