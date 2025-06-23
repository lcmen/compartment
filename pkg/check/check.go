package check

import (
	"compartment/pkg/container"
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"net"
	"time"
)

type result struct {
	valid bool
	err   error
}

func Check() error {
	res := &result{valid: true}

	checkDocker(res)
	checkDevDNS(res)
	checkDirectAccess(res)

	if res.valid {
		fmt.Println("All checks passed successfully.")
	}

	return res.err
}

func checkDocker(res *result) {
	if !res.valid {
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		res.err = fmt.Errorf("failed to create Docker client: %w", err)
		res.valid = false
		return
	}

	defer cli.Close()

	_, err = cli.Ping(context.Background())

	if err != nil {
		fmt.Println("Docker is not running or is unreachable")
		res.valid = false
	}
}

func checkDevDNS(res *result) {
	if !res.valid {
		return
	}

	cnt, err := container.NewContainer("devdns")
	if err != nil {
		res.err = fmt.Errorf("failed to create devdns container: %w", err)
		res.valid = false
		return
	}
	defer cnt.Close()

	if cnt.State != container.StateRunning {
		fmt.Println("DevDNS container is not running. To start it, run: compartment start devdns")
		res.valid = false
	}
}

func checkDirectAccess(res *result) {
	if !res.valid {
		return
	}

	container, err := container.ExistingContainer("devdns")
	if err != nil {
		res.err = fmt.Errorf("failed to get devdns container: %w", err)
		res.valid = false
		return
	}

	ip, err := container.IPAddress()
	if err != nil {
		res.err = fmt.Errorf("failed to get direct access IP: %w", err)
		res.valid = false
		return
	}

	address := net.JoinHostPort(ip, "53")
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Println("Could not connect to Docker container via its IP. If you are on MacOS make sure that `docker-mac-net-connect` is running.")
		res.valid = false
		return
	}

	conn.Close()
}
