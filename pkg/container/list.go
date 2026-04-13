package container

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func List() ([]Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	f := filters.NewArgs(filters.Arg("label", "compartment=true"))

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: f,
	})
	if err != nil {
		return nil, err
	}

	result := make([]Container, 0, len(containers))
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		var state State
		switch c.State {
		case "running":
			state = StateRunning
		case "exited":
			state = StateStopped
		default:
			state = StateRemoved
		}

		result = append(result, Container{
			Name:   name,
			Image:  c.Image,
			State:  state,
			Labels: c.Labels,
		})
	}

	return result, nil
}
