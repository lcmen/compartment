package service

import (
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

var defaultDevDNSEnv = []string{
	"DNS_DOMAIN=container",
}

func NewDevDNSService(env []string) (*Service, error) {
	env = append(defaultDevDNSEnv, env...)
	volumes := getDevDNSVolumes()
	ports := getDevDNSPorts()

	return &Service{
		Name:    "devdns",
		Image:   "ruudud/devdns:latest",
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

func getDevDNSPorts() nat.PortMap {
	port := nat.Port("53/udp")
	return nat.PortMap{
		port: []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "53",
			},
		},
	}
}
