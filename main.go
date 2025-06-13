package main

import (
	"compartment/cmd"
	"compartment/pkg/logging"
	"os"
)

func main() {
	if err := cmd.Run(); err != nil {
		logging.Error(err.Error())
		os.Exit(1)
	}
}
