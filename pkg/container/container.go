package container

import (
	"compartment/pkg/logging"
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/go-connections/nat"
)

type State int

const (
	StateRunning State = iota
	StateStopped
	StateRemoved
	StateError
)

type Container struct {
	Name  string
	Image string
	State State
	Err   error
	cid   string
	cli   *client.Client
}

func NewContainer(name string) (*Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containerJSON, err := cli.ContainerInspect(context.Background(), name)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return nil, err
		} else {
			return &Container{Name: name, State: StateRemoved, cli: cli}, nil
		}
	}

	var cid string
	var state State = StateRemoved
	if containerJSON.State.Running {
		cid = containerJSON.ID
		state = StateRunning
	} else if containerJSON.State.Status == "exited" {
		cid = containerJSON.ID
		state = StateStopped
	} else {
		state = StateRemoved
	}

	return &Container{Name: name, State: state, cid: cid, cli: cli}, nil
}

func ExistingContainer(name string) (*Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containerJSON, err := cli.ContainerInspect(context.Background(), name)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return nil, err
		} else {
			return &Container{Name: name, State: StateRemoved, cli: cli}, nil
		}
	}

	return &Container{
		Name:  name,
		Image: containerJSON.Image,
		State: StateRunning,
		cid:   containerJSON.ID,
		cli:   cli,
	}, nil
}

func (c *Container) Create(image string, env []string, volumes []mount.Mount, ports nat.PortMap) error {
	if c.State != StateRemoved {
		return fmt.Errorf("container not in removable state: %v", c.State)
	}

	if err := PullImage(c.cli, context.Background(), image); err != nil {
		logging.Debug(fmt.Sprintf("error pulling image %s: %v", image, err))
		c.State = StateError
		c.Err = err
		return err
	}

	resp, err := c.cli.ContainerCreate(
		context.Background(),
		&container.Config{Image: image, Env: env, ExposedPorts: exposedPorts(ports)},
		&container.HostConfig{AutoRemove: true, Mounts: volumes, PortBindings: ports},
		nil,
		nil,
		c.Name,
	)
	if err != nil {
		logging.Debug(fmt.Sprintf("error creating container %s: %v", c.Name, err))
		c.State = StateError
		c.Err = err
		return err
	}

	c.cid = resp.ID
	c.State = StateStopped
	return nil
}

func (c *Container) Close() {
	c.cli.Close()
}

func (c *Container) Stop() {
	if c.State != StateRunning {
		return
	}

	logging.Debug(fmt.Sprintf("stopping container %s", c.Name))
	err := c.cli.ContainerStop(context.Background(), c.Name, container.StopOptions{})
	if err != nil {
		logging.Debug(fmt.Sprintf("error stopping container %s: %v", c.Name, err))
		c.State = StateError
		c.Err = err
		return
	}

	c.State = StateStopped
}

func (c *Container) Remove() {
	if c.State != StateStopped {
		return
	}

	logging.Debug(fmt.Sprintf("removing container %s", c.Name))
	err := c.cli.ContainerRemove(context.Background(), c.Name, container.RemoveOptions{Force: true})
	if err != nil {
		logging.Debug(fmt.Sprintf("error removing container %s: %v", c.Name, err))
		c.State = StateError
		c.Err = err
		return
	}

	c.State = StateRemoved
}

func (c *Container) Start() {
	if c.State != StateStopped {
		return
	}

	logging.Debug(fmt.Sprintf("starting container %s", c.Name))
	err := c.cli.ContainerStart(context.Background(), c.cid, container.StartOptions{})
	if err != nil {
		logging.Debug(fmt.Sprintf("error starting container %s: %v", c.Name, err))
		c.State = StateError
		c.Err = err
		return
	}

	c.State = StateRunning
}

func exposedPorts(ports nat.PortMap) nat.PortSet {
	exposedPorts := make(nat.PortSet)
	for port := range ports {
		exposedPorts[port] = struct{}{}
	}
	return exposedPorts
}
