package container

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func PullImage(cli *client.Client, ctx context.Context, name string) error {
	fmt.Printf("Updating image %s...\n", name)
	reader, err := cli.ImagePull(ctx, name, image.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	dec := json.NewDecoder(reader)
	var msg struct {
		Error string `json:"error"`
	}
	for {
		if err := dec.Decode(&msg); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if msg.Error != "" {
			return fmt.Errorf("error pulling image: %s", msg.Error)
		}
	}
	fmt.Println("Image downloaded.")
	return nil
}
