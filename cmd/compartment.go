package cmd

import (
	"flag"
	"fmt"
	"os"
)

type Service struct {
	Name string
	Image string
	Version string
	Env []string
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `
   ___                          _                 _
  / __|___ _ __  _ __  __ _ _ _| |_ _ __  ___ _ _| |_
 | (__/ _ \ '  \| '_ \/ _| | '_|  _| '  \/ -_) ' \  _|
  \___\___/_|_|_| .__/\__,_|_|  \__|_|_|_\___|_||_\__|
                |_|
Compartment - your local service runner

Usage:
  compartment [flags] <command> <service> [version]

Flags:
`)
		flag.PrintDefaults()
	}
}

func Run() {
	service, command, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch command {
	case "start":
		fmt.Printf("Starting service: %s\n", service.Image)
		fmt.Printf("Container name: %s\n", service.Name)
		fmt.Printf("Environment variables: %v\n", service.Env)
		// TODO: Implement Docker logic here
	case "stop":
		fmt.Printf("Stopping service: %s\n", service.Image)
		fmt.Printf("Container name: %s\n", service.Name)
		// TODO: Implement Docker logic here
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: start, stop")
		os.Exit(1)
	}
}

func parseArgs() (Service, string, error) {
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

	if len(args) < 2 {
		return Service{}, "", fmt.Errorf("Usage: compartment [flags] <command> <service> [version]")
	}

	command, service, version := getArgOrDefault(args, 0, ""), getArgOrDefault(args, 1, ""), getArgOrDefault(args, 2, "latest")
	if name == "" && version == "latest" {
		name = service
	} else if name == "" {
		name = fmt.Sprintf("%s%s", service, version)
	}

	return Service{
		Name: name,
		Image: fmt.Sprintf("%s:%s", service, version),
		Env: envs,
	}, command, nil
}

func getArgOrDefault(args []string, i int, def string) string {
	if i < len(args) {
		return args[i]
	}
	return def
}
