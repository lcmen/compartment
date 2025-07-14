package cmd

import (
	"compartment/pkg/check"
	"compartment/pkg/service"
	"flag"
	"fmt"
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

Commands:
  start     Start a service
  stop      Stop a service
  status    Show service status
  check     Check devdns container
  help      Show this help message

Services:
  postgres  PostgreSQL database
  redis     Redis cache
  devdns    Development DNS server

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
		return err
	}

	switch command {
	case "start":
		err := service.Start()
		if err != nil {
			return fmt.Errorf("error starting container: %w", err)
		}
	case "stop":
		err := service.Stop()
		if err != nil {
			return fmt.Errorf("error stopping container: %w", err)
		}
	case "status":
		err := service.Status()
		if err != nil {
			return fmt.Errorf("error getting container status: %w", err)
		}
	case "check":
		err := check.Check()
		if err != nil {
			return fmt.Errorf("error checking devdns container: %w", err)
		}
	case "help":
		flag.Usage()
		return nil
	default:
		return fmt.Errorf("unknown command: %s (available commands: start, stop, status, check)", command)
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
	srv, err := getServiceForCommand(cmd, name, args)

	return srv, cmd, err
}

func getServiceForCommand(cmd, name string, args []string) (*service.Service, error) {
	if cmd == "check" || cmd == "help" {
		return nil, nil
	}

	if len(args) < 2 {
		return nil, fmt.Errorf("command '%s' requires at least <service> argument", cmd)
	}

	kind := getArgOrDefault(args, 1, "")
	ver := getArgOrDefault(args, 2, "latest")
	serviceName := getServiceName(kind, ver)

	containerName := serviceName
	if name != "" {
		containerName = name
	}

	if cmd == "stop" || cmd == "status" {
		return service.NewServiceFromExistingContainer(containerName)
	} else {
		return service.NewService(containerName, kind, ver, nil)
	}
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
