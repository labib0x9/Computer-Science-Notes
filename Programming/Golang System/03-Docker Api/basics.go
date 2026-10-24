package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func main() {

	imageN := "alpine:3.20"
	containerId := ""

	// --- // Connect to dockerhub
	client, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Pull image if not present...
	imgResp, err := client.ImageInspect(context.Background(), imageN)
	if err != nil {
		if strings.Contains(err.Error(), "No such image:") {
			reader, err := client.ImagePull(context.Background(), imageN, image.PullOptions{})
			if err != nil {
				panic(err)
			}
			io.Copy(os.Stdout, reader)
			reader.Close()
		} else {
			panic(err)
		}
	}
	_ = imgResp

	// client.ImagePull(context.Background(), "")

	containerN := "container-vv9"
	resp, err := client.ContainerCreate(context.Background(),
		&container.Config{
			Image: imageN, // image has to be present.. or pull
			Cmd:   []string{"echo", "Hello host"},
		},
		nil, nil, nil, containerN,
	)
	if err != nil {
		panic(err)
	}

	containerId = resp.ID
	fmt.Println("Continaer", containerN, "Created# Id:", containerId)

	// --- // STart the container.
	if err := client.ContainerStart(context.Background(), containerId, container.StartOptions{}); err != nil {
		panic(err)
	}

	// --- // Wait container to finish, Not Running -> Blocks until the container exits
	respCh, errCh := client.ContainerWait(context.Background(), containerId, container.WaitConditionNotRunning)
	select {
	case resp := <-respCh:
		fmt.Println(resp.StatusCode)
	case err := <-errCh:
		fmt.Println(err.Error())
	}

	// --- // Container Logs
	logs, err := client.ContainerLogs(context.Background(), containerId, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		panic(err)
	}
	defer logs.Close()

	io.Copy(os.Stdout, logs)

	// --- // Remove Container
	if err := client.ContainerRemove(context.Background(), containerId, container.RemoveOptions{}); err != nil {
		panic(err)
	}
}
