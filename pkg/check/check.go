package check

import (
	"compartment/pkg/container"
	"fmt"
)

func Check() error {
	return CheckDevDNS()
}

func CheckDevDNS() error {
	cnt, err := container.NewContainer("devdns")
	if err != nil {
		return err
	}
	defer cnt.Close()

	if cnt.State == container.StateRunning {
		fmt.Println("devdns container is running.")
		return nil
	}

	fmt.Println("devdns is not running.")
	fmt.Println("To start it, run:")
	fmt.Println("compartment start devdns")
	return nil
}
