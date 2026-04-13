package cmd

import (
	"compartment/pkg/container"
	"fmt"
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

	fmt.Printf("%-20s %-12s %-10s %s\n", "NAME", "SERVICE", "VERSION", "STATE")
	for _, c := range containers {
		service := c.Labels["compartment.service"]
		version := c.Labels["compartment.version"]

		var state string
		switch c.State {
		case container.StateRunning:
			state = "running"
		case container.StateStopped:
			state = "stopped"
		default:
			state = "removed"
		}

		fmt.Printf("%-20s %-12s %-10s %s\n", c.Name, service, version, state)
	}

	return nil
}
