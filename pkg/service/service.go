package service

import (
	"fmt"
	"compartment/pkg/container"
	"compartment/pkg/logging"
	"github.com/docker/docker/api/types/mount"
)

type Service struct {
	Name string
	Image string
	Version string
	Env []string
	Volumes []mount.Mount
}

func NewService(name, service, version string, env []string) (*Service, error) {
	if name == "" {
		return nil, fmt.Errorf("service name cannot be empty")
	}
	if service == "" {
			return nil, fmt.Errorf("service type cannot be empty")
	}

	logging.Debug(fmt.Sprintf("Creating %s service called %s", service, name))

	switch service {
	case "postgres":
		return NewPostgresService(name, service, version, env)
	case "redis":
		return NewRedisService(name, service, version, env)
	default:
		return nil, fmt.Errorf("Unknown service: %s, only postgres and redis are supported", service)
	}
}

func (s *Service) Start() error {
	cnt := container.NewContainer(s.Name)

	if cnt.State == container.StateRunning {
		printServiceInfo(cnt, s)
		return nil
	}

	if cnt.State == container.StateStopped {
		logging.Debug(fmt.Sprintf("Removing stopped container %s to start it again", cnt.Name))
		cnt.Remove()
	}

	cnt.Create(s.Image, s.Env, s.Volumes)
	cnt.Start()

	printServiceInfo(cnt, s)

	return cnt.Err
}

func (s *Service) Status() error {
	cnt := container.NewContainer(s.Name)

	printServiceInfo(cnt, s)

	return cnt.Err
}

func (s *Service) Stop() error {
	cnt := container.NewContainer(s.Name)

	cnt.Stop()

	printServiceInfo(cnt, s)

	return cnt.Err
}

func printServiceInfo(cnt *container.Container, s *Service) {
	switch cnt.State {
	case container.StateRunning:
		fmt.Printf("Service `%s` using `%s` image is running.\n", cnt.Name, s.Image)
	case container.StateStopped:
		fmt.Printf("Service `%s` using `%s` image has stopped.\n", cnt.Name, s.Image)
	case container.StateRemoved:
		fmt.Printf("Service `%s` using `%s` image is not running.\n", cnt.Name, s.Image)
	}
}
