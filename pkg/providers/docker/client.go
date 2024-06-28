/*
Copyright © 2024 Alessandro Sanino <alessandro@sanino.dev>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

// Package docker contains the Docker wrapped client for the cat image.
package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"

	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type onPullProgressFunc func(reader io.Reader)

// DockerClient is a wrapper around the Docker client
// to handle the CLI features.
type DockerClient struct {
	docker         *docker.Client
	onPullProgress onPullProgressFunc
}

// NewDockerClient creates a new DockerClient
func NewDockerClient(onPullProgress onPullProgressFunc) (*DockerClient, error) {
	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &DockerClient{
		docker:         dockerClient,
		onPullProgress: onPullProgress,
	}, nil
}

// Close closes the underlying Docker client
func (client *DockerClient) Close() error {
	return client.docker.Close()
}

// PullCatImage pulls the specified version of the cat image
func (client *DockerClient) PullCatImage(ctx context.Context, version string) error {
	fullImageURL := fmt.Sprintf("cheshire-cat-ai:%s", version)

	slog.Debug("Pulling image", slog.String("image", fullImageURL))
	result, err := client.docker.ImagePull(ctx, fullImageURL, image.PullOptions{})
	if err != nil {
		return err
	}
	defer result.Close()

	if client.onPullProgress == nil {
		io.Copy(io.Discard, result)
	} else {
		proxyBuffer := new(bytes.Buffer)
		client.onPullProgress(proxyBuffer)
		io.Copy(proxyBuffer, result)
	}

	return nil
}

// https://docs.docker.com/engine/api/sdk/examples/

type StartCatContainerConfig struct {
	CatImage              string
	CatContainerName      string
	CatContainerBoundPort int
	PluginFolderPath      string
	DataFolderPath        string
	StaticFolderPath      string
}

// StartCatContainer starts the cheshire cat container with the specified config
func (client *DockerClient) StartCatContainer(ctx context.Context, config StartCatContainerConfig) error {
	dockerContainerConfig := &container.Config{
		Tty:   false,
		Image: config.CatImage,
		Volumes: map[string]struct{}{
			config.PluginFolderPath: {},
		},
	}

	portBinding := nat.Port(fmt.Sprintf("%d/tcp", config.CatContainerBoundPort))
	dockerHostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:%s", config.PluginFolderPath, "/app/cat/plugins"),
			fmt.Sprintf("%s:%s", config.DataFolderPath, "/app/cat/data"),
			fmt.Sprintf("%s:%s", config.StaticFolderPath, "/app/cat/static"),
		},
		PortBindings: map[nat.Port][]nat.PortBinding{
			portBinding: {{HostPort: "80"}},
		},
	}
	result, err := client.docker.ContainerCreate(ctx, dockerContainerConfig, dockerHostConfig, nil, nil, config.CatContainerName)
	if err != nil {
		return err
	}

	err = client.docker.ContainerStart(ctx, result.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	statusCh, errCh := client.docker.ContainerWait(ctx, result.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	return nil
}

// StopCatContainer stops the specified cat container
func (client *DockerClient) StopCatContainer(ctx context.Context, containerName string) error {
	return client.docker.ContainerStop(ctx, containerName, container.StopOptions{})
}

// RemoveCatContainer removes the specified cat container
func (client *DockerClient) RemoveCatContainer(ctx context.Context, containerName string) error {
	return client.docker.ContainerRemove(ctx, containerName, container.RemoveOptions{
		Force:         true,
		RemoveVolumes: false,
	})
}
