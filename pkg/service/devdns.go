package service

import (
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"strings"
)

var defaultDevDNSEnv = []string{
	"DNS_DOMAIN=container",
}

func NewDevDNSService(env []string) (*Service, error) {
	env = append(defaultDevDNSEnv, env...)
	volumes := getDevDNSVolumes()
	ports := getDevDNSPorts(envValue(env, "DNS_PORT", "53"))

	return &Service{
		Name:    "devdns",
		Image:   "lmendelowski/devdns:latest",
		Kind:    "devdns",
		Version: "latest",
		Env:     env,
		Volumes: volumes,
		Ports:   ports,
	}, nil
}

func getDevDNSVolumes() []mount.Mount {
	return []mount.Mount{
		{
			Type:     mount.TypeBind,
			Source:   "/var/run/docker.sock",
			Target:   "/var/run/docker.sock",
			ReadOnly: true,
		},
	}
}

func getDevDNSPorts(hostPort string) nat.PortMap {
	port := nat.Port("53/udp")
	return nat.PortMap{
		port: []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: hostPort,
			},
		},
	}
}

func envValue(env []string, key, def string) string {
	prefix := key + "="
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			return strings.TrimPrefix(e, prefix)
		}
	}
	return def
}
