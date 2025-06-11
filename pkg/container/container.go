package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
)

type State int
const (
	StateRunning State = iota
	StateStopped
	StateRemoved
	StateError
)

type Container struct {
	Name string
	Image string
	State State
	Err error
	cid string
	cli *client.Client
}

func NewContainer(name string) *Container {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return &Container{Name: name, State: StateError, Err: err}
	}
	defer cli.Close()

	containerJSON, err := cli.ContainerInspect(context.Background(), name)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return &Container{Name: name, State: StateError, Err: err}
		} else {
			return &Container{Name: name, State: StateRemoved, cli: cli}
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

	return &Container{Name: name, State: state, cid: cid, cli: cli}
}

func (c *Container) Create(image string, env []string, volumes []mount.Mount) error {
	if c.State != StateRemoved {
		return fmt.Errorf("container not in removable state: %v", c.State)
	}

	if err := PullImage(c.cli, context.Background(), image); err != nil {
		c.State = StateError
		c.Err = err
		return err
	}

	resp, err := c.cli.ContainerCreate(
		context.Background(),
		&container.Config{Image: image, Env: env},
		&container.HostConfig{AutoRemove: true, Mounts: volumes},
		nil,
		nil,
		c.Name,
	)
	if err != nil {
		c.State = StateError
		c.Err = err
		return err
	}

	c.cid = resp.ID
	c.State = StateStopped
	return nil
}

func (c *Container) Stop() {
	if c.State != StateRunning {
		return
	}

	err := c.cli.ContainerStop(context.Background(), c.Name, container.StopOptions{})
	if err != nil {
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

	err := c.cli.ContainerRemove(context.Background(), c.Name, container.RemoveOptions{Force: true})
	if err != nil {
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

	err := c.cli.ContainerStart(context.Background(), c.cid, container.StartOptions{})
	if err != nil {
		c.State = StateError
		c.Err = err
		return
	}

	c.State = StateRunning
}
