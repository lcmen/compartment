package service

import (
	"compartment/pkg/configuration"
	"fmt"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"os"
)

const defaultRedisVersion = "8"

var defaultRedisEnv = []string{}

func NewRedisService(name, service, version string, env []string) (*Service, error) {
	if version == "" {
		version = defaultRedisVersion
	}
	if name == "" {
		name = fmt.Sprintf("%s.%s", version, service)
	}
	env = append(defaultRedisEnv, env...)
	volumes, err := getRedisVolumes(name)
	if err != nil {
		return nil, err
	}

	return &Service{
		Name:    name,
		Image:   fmt.Sprintf("%s:%s", service, version),
		Version: version,
		Env:     env,
		Volumes: volumes,
		Ports:   make(nat.PortMap),
	}, nil
}

func getRedisVolumes(name string) ([]mount.Mount, error) {
	source, err := prepareRedisDataDir(name)
	if err != nil {
		return nil, err
	}

	return []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: source,
			Target: "/data",
		},
	}, nil
}

func prepareRedisDataDir(name string) (string, error) {
	dataDir := configuration.Get().DataDir

	dir := fmt.Sprintf("%s/redis/%s", dataDir, name)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	return dir, nil
}
