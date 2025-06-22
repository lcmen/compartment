package service

import (
	"compartment/pkg/configuration"
	"fmt"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"os"
)

var defaultPostgresEnv = []string{
	"POSTGRES_USER=postgres",
	"POSTGRES_PASSWORD=postgres",
	"POSTGRES_DB=postgres",
}

func NewPostgresService(name, service, version string, env []string) (*Service, error) {
	env = append(defaultPostgresEnv, env...)
	volumes := getPostgresVolumes(name)

	return &Service{
		Name:    name,
		Image:   fmt.Sprintf("%s:%s", service, version),
		Version: version,
		Env:     env,
		Volumes: volumes,
		Ports:   make(nat.PortMap),
	}, nil
}

func getPostgresVolumes(name string) []mount.Mount {
	return []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: preparePostgresDataDir(name),
			Target: "/var/lib/postgresql/data",
		},
	}
}

func preparePostgresDataDir(name string) string {
	dataDir := configuration.Get().DataDir
	dir := fmt.Sprintf("%s/postgres/%s", dataDir, name)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	return dir
}
