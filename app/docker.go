package app

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"strings"
	"time"
)

func WaitForContainer(containerName string) error {
	// Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error initializing docker client: %w", err)
	}
	// KubeContext to set a timeout
	timeoutContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for !ContainerExists(containerName, dockerClient, timeoutContext) {
	}
	for {
		// Get container details
		containerJSON, err := dockerClient.ContainerInspect(timeoutContext, containerName)
		if err != nil {
			return fmt.Errorf("error inspecting container: %w", err)
		}

		// Check if the container is running
		if containerJSON.State.Running {
			return nil
		}

		// Wait for a short interval before checking again
		time.Sleep(1 * time.Second)
	}
}

func ContainerExists(containerName string, dockerClient *client.Client, timeoutContext context.Context) bool {
	_, err := dockerClient.ContainerInspect(timeoutContext, containerName)
	if err != nil && strings.Contains(err.Error(), "No such container") {
		return false
	}
	return true
}
