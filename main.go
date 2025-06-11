package main

import (
	"compartment/cmd"
	"fmt"
	"os"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
  }
}
