package container

import (
	"context"
	"encoding/json"
	"compartment/pkg/logging"
	"fmt"
	"io"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type ImageMessage struct {
	Status   string `json:"status"`
	Progress string `json:"progress"`
	ID       string `json:"id"`
	Error    string `json:"error"`
}

func PullImage(cli *client.Client, ctx context.Context, name string) error {
	logging.Info(fmt.Sprintf("Updating image %s...", name))
	reader, err := cli.ImagePull(ctx, name, image.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	dec := json.NewDecoder(reader)
	var msg ImageMessage
	lastStatus := make(map[string]string)
	for {
		if err := dec.Decode(&msg); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if msg.Error != "" {
			return fmt.Errorf("error pulling image: %s", msg.Error)
		}
		key := msgKey(&msg)
		if lastStatus[key] != msg.Status {
			logMsg(&msg)
			lastStatus[key] = msg.Status
		}
	}
	return nil
}

func logMsg(msg *ImageMessage) {
	if msg.ID != "" {
		logging.Info(fmt.Sprintf("%s: %s", msg.ID, msg.Status))
	} else {
		logging.Info(msg.Status)
	}
}

func msgKey(msg *ImageMessage) string {
	if msg.ID != "" {
		return msg.ID
	}
	return "_global"
}
