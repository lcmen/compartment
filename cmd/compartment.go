package cmd

import (
	"flag"
	"fmt"
	"os"
	"compartment/pkg/service"
)

const help = "Usage: compartment [flags] <command> <service> [version]\n\nUse `--help` command to display available options."
const usage = `
   ___                          _                 _
  / __|___ _ __  _ __  __ _ _ _| |_ _ __  ___ _ _| |_
 | (__/ _ \ '  \| '_ \/ _| | '_|  _| '  \/ -_) ' \  _|
  \___\___/_|_|_| .__/\__,_|_|  \__|_|_|_\___|_||_\__|
                |_|

Compartment - your assistant for spinning up services needed for local development in a Docker-based environment.

Usage:
  compartment [flags] <command> <service> [version]

Flags:
`

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage)
		flag.PrintDefaults()
	}
}

func Run() error {
	service, command, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch command {
	case "start":
		err := service.Start()
		if err != nil {
			return fmt.Errorf("Error starting container: %v\n", err)
		}
	case "stop":
		err := service.Stop()
		if err != nil {
			return fmt.Errorf("Error stopping container: %v\n", err)
		}
	case "status":
		err := service.Status()
		if err != nil {
			return fmt.Errorf("Error getting container status: %v\n", err)
		}
	case "help":
		flag.Usage()
		return nil
	default:
		return fmt.Errorf("Unknown command: %s\nAvailable commands: start, stop", command)
	}

	return nil
}

func parseArgs() (*service.Service, string, error) {
	args := []string{}
	envs := []string{}
	name := ""

	flag.StringVar(&name, "n", "", "Custom container name")
	flag.Func("e", "Environment variable (can be specified multiple times)", func(val string) error {
		envs = append(envs, val)
		return nil
	})
	flag.Parse()

	args = flag.Args()

	if len(args) == 0 {
		return nil, "help", nil
	}

	cmd := getArgOrDefault(args, 0, "")

	// If command is "stop" and -n flag is provided, allow skipping service and version
	if cmd == "stop" && name != "" {
		srv, err := service.NewServiceFromExistingContainer(name)
		if err != nil {
			return nil, "", err
		}
		return srv, cmd, nil
	}

	// For all other cases, require at least <command> and <service>
	if len(args) < 2 {
		return nil, "", fmt.Errorf(help)
	}

	kind := getArgOrDefault(args, 1, "")
	ver := getArgOrDefault(args, 2, "latest")
	serviceName := getServiceName(kind, ver)

	srv, err := service.NewService(serviceName, kind, ver, envs)
	if err != nil {
		return nil, "", err
	}

	return srv, cmd, nil
}

func getServiceName(kind string, ver string) string {
	if ver == "latest" {
		return kind
	}
	return fmt.Sprintf("%s.%s", ver, kind)
}

func getArgOrDefault(args []string, i int, def string) string {
	if i < len(args) {
		return args[i]
	}
	return def
}
