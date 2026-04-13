package service

import (
	"compartment/pkg/configuration"
	"fmt"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"os"
	"strconv"
)

const defaultPostgresVersion = "18"

var defaultPostgresEnv = []string{
	"POSTGRES_USER=postgres",
	"POSTGRES_PASSWORD=postgres",
	"POSTGRES_DB=postgres",
}

func NewPostgresService(name, service, version string, env []string) (*Service, error) {
	if version == "" {
		version = defaultPostgresVersion
	}
	if name == "" {
		name = fmt.Sprintf("%s.%s", version, service)
	}
	env = append(defaultPostgresEnv, env...)
	volumes, err := getPostgresVolumes(name, version)
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

func getPostgresVolumes(name, version string) ([]mount.Mount, error) {
	var target string
	major, err := strconv.Atoi(version)
	if err != nil || major >= 18 {
		target = "/var/lib/postgresql"
	} else {
		target = "/var/lib/postgresql/data"
	}

	source, err := preparePostgresDataDir(name)
	if err != nil {
		return nil, err
	}

	return []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: source,
			Target: target,
		},
	}, nil
}

func preparePostgresDataDir(name string) (string, error) {
	dataDir := configuration.Get().DataDir
	dir := fmt.Sprintf("%s/postgres/%s", dataDir, name)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	return dir, nil
}
