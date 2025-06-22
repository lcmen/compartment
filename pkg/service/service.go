package service

import (
	"compartment/pkg/container"
	"compartment/pkg/logging"
	"fmt"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

type Service struct {
	Name    string
	Image   string
	Version string
	Env     []string
	Volumes []mount.Mount
	Ports   nat.PortMap
}

func NewService(name, service, version string, env []string) (*Service, error) {
	if name == "" {
		return nil, fmt.Errorf("service name cannot be empty")
	}
	if service == "" {
		return nil, fmt.Errorf("service type cannot be empty")
	}

	logging.Debug(fmt.Sprintf("Initializing %s service called %s", service, name))

	switch service {
	case "postgres":
		return NewPostgresService(name, service, version, env)
	case "redis":
		return NewRedisService(name, service, version, env)
	case "devdns":
		return NewDevDNSService(env)
	default:
		return nil, fmt.Errorf("unknown service: %s (only postgres, redis, and devdns are supported)", service)
	}
}

func NewServiceFromExistingContainer(name string) (*Service, error) {
	logging.Debug(fmt.Sprintf("Deriving service from existing container %s", name))

	cnt, err := container.ExistingContainer(name)
	if err != nil {
		return nil, err
	}
	defer cnt.Close()

	if cnt.State == container.StateError {
		return nil, cnt.Err
	}

	return &Service{Name: name, Image: cnt.Image}, nil
}

func (s *Service) Start() error {
	cnt, err := container.NewContainer(s.Name)
	if err != nil {
		return err
	}
	defer cnt.Close()

	if cnt.State == container.StateRunning {
		printServiceInfo(cnt, s)
		return nil
	}

	if cnt.State == container.StateStopped {
		logging.Debug(fmt.Sprintf("Removing stopped container %s to start it again", cnt.Name))
		cnt.Remove()
	}

	cnt.Create(s.Image, s.Env, s.Volumes, s.Ports)
	cnt.Start()

	printServiceInfo(cnt, s)

	return cnt.Err
}

func (s *Service) Status() error {
	cnt, err := container.NewContainer(s.Name)
	if err != nil {
		return err
	}
	defer cnt.Close()

	printServiceInfo(cnt, s)

	return cnt.Err
}

func (s *Service) Stop() error {
	cnt, err := container.NewContainer(s.Name)
	if err != nil {
		return err
	}
	defer cnt.Close()

	cnt.Stop()

	printServiceInfo(cnt, s)

	return cnt.Err
}

func printServiceInfo(cnt *container.Container, s *Service) {
	imageInfo := ""
	if cnt.Image != "" {
		imageInfo = fmt.Sprintf(" using `%s` image", cnt.Image)
	}

	switch cnt.State {
	case container.StateRunning:
		fmt.Printf("Service `%s`%s is running.\n", cnt.Name, imageInfo)
	case container.StateStopped:
		fmt.Printf("Service `%s`%s has stopped.\n", cnt.Name, imageInfo)
	case container.StateRemoved:
		fmt.Printf("Service `%s`%s is not running.\n", cnt.Name, imageInfo)
	}
}
