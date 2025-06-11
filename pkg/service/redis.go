package service

import (
	"fmt"
	"os"
	"github.com/docker/docker/api/types/mount"
	"compartment/pkg/configuration"
)

var defaultRedisEnv = []string{}

func NewRedisService(name, service, version string, env []string) (*Service, error) {
	env = append(defaultRedisEnv, env...)
	volumes := getRedisVolumes(name)

	return &Service{
		Name: name,
		Image: fmt.Sprintf("%s:%s", service, version),
		Version: version,
		Env: env,
		Volumes: volumes,
	}, nil
}

func getRedisVolumes(name string) []mount.Mount {
	return []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: prepareRedisDataDir(name),
			Target: "/data",
		},
	}
}

func prepareRedisDataDir(name string) string {
	dataDir := configuration.Get().DataDir

	dir := fmt.Sprintf("%s/redis/%s", dataDir, name)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	return dir
}
