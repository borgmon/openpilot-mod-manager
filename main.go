/*
Copyright © 2022 borgmon

*/
package main

import (
	"os"

	"github.com/borgmon/openpilot-mod-manager/cmd"
)

func main() {
	cmd.Execute()
}

func GetEnvWithDefault(env string, def string) string {
	if env, ok := os.LookupEnv(env); ok {
		return env
	} else {
		return def
	}
}
