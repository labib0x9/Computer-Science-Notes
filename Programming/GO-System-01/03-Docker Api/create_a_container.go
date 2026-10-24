package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {

	// --- // Connect to dockerhub
	client, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// --- // Create Container
	containerN := "container-vv2"
	resp, err := client.ContainerCreate(context.Background(),
		&container.Config{
			Image: "alpine:3.20", // image has to be present.. or pull
			Cmd:   []string{"echo", "Hello host"},
		},
		nil, nil, nil, containerN,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Continaer", containerN, "Created# Id:", resp.ID)

}
