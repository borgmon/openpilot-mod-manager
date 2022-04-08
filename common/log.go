package common

import (
	"fmt"

	"github.com/pkg/errors"
)

func PanicIfErr(err error) error {
	if err != nil {
		panic(errors.WithStack(err))
	}
	return err
}

func LogIfErr(err error) error {
	if err != nil {
		fmt.Printf("%+v", err)
	}
	return err
}
