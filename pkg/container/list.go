package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func List() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	f := filters.NewArgs(filters.Arg("label", "compartment=true"))
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{Filters: f})
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		fmt.Println("No compartment-managed containers found.")
		return nil
	}

	fmt.Printf("%-16s %-24s %-12s %s\n", "NAME", "CONTAINER", "SERVICE", "VERSION")
	for _, c := range containers {
		name := strings.TrimPrefix(c.Names[0], "/")
		svc := c.Labels["compartment.service"]
		ver := c.Labels["compartment.version"]
		fmt.Printf("%-16s %-24s %-12s %s\n", strings.TrimSuffix(name, "."+svc), name, svc, ver)
	}

	return nil
}
