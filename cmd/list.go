package cmd

import (
	"compartment/pkg/container"
	"fmt"
	"strings"
)

func listCmd() error {
	containers, err := container.List()
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		fmt.Println("No compartment-managed containers found.")
		return nil
	}

	fmt.Printf("%-16s %-24s %-12s %s\n", "NAME", "CONTAINER", "SERVICE", "VERSION")
	for _, c := range containers {
		service := c.Labels["compartment.service"]
		version := c.Labels["compartment.version"]
		name := strings.TrimSuffix(c.Name, "."+service)

		fmt.Printf("%-16s %-24s %-12s %s\n", name, c.Name, service, version)
	}

	return nil
}
