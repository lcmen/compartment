package configuration

import (
	"fmt"
	"log/slog"
	"os"
)

type Configuration struct {
	DataDir string
	Logger  *slog.Logger
}

var config *Configuration

func init() {
	var loglevel slog.Level
	dataDir := fmt.Sprintf("%s/compartment", os.Getenv("XDG_DATA_HOME"))

	if os.Getenv("DEBUG") == "1" {
		loglevel = slog.LevelDebug
	} else {
		loglevel = slog.LevelInfo
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			panic(fmt.Sprintf("failed to create data directory %s: %v", dataDir, err))
		}
	}

	config = &Configuration{
		DataDir: dataDir,
		Logger:  slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: loglevel})),
	}
}

func Get() *Configuration {
	return config
}
