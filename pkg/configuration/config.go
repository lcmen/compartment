package configuration

import (
	"fmt"
	"os"
)

type Configuration struct {
	DataDir string
}

var config *Configuration

func init() {
	dataDir := fmt.Sprintf("%s/compartment", os.Getenv("XDG_DATA_HOME"))

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0755)
	}

	config = &Configuration{
		DataDir: dataDir,
	}
}

func Get() *Configuration {
	return config
}
