package models

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	client *client.Client
	ctx    context.Context
}

// NewDockerClient creates a new Docker client
// Returns a pointer to a new DockerClient
// Panics on error
func NewDockerClient() *DockerClient {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &DockerClient{
		client: client,
		ctx:    context.Background(),
	}
}

// Close closes the Docker client
// Returns an error if the client cannot be closed
func (d *DockerClient) Close() error {
	return d.client.Close()
}

// GetDiskUsage returns the disk usage of the Docker client
// Returns an error if the disk usage cannot be retrieved
func (d *DockerClient) GetDiskUsage() (*types.DiskUsage, error) {
	usage, err := d.client.DiskUsage(d.ctx, types.DiskUsageOptions{})
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

// GetStoppedContainers returns a list of stopped containers
// Returns an error if the list cannot be retrieved
func (d *DockerClient) GetStoppedContainers() ([]types.Container, error) {
	args := filters.NewArgs()
	args.Add("status", "exited")
	args.Add("status", "created")
	args.Add("status", "dead")

	return d.client.ContainerList(d.ctx, container.ListOptions{All: true, Filters: args})
}

// RemoveContainer removes a container
// Returns an error if the container cannot be removed
func (d *DockerClient) RemoveContainer(containerID string) error {
	return d.client.ContainerRemove(d.ctx, containerID, container.RemoveOptions{
		RemoveVolumes: false,
		Force:         false,
	})
}

// GetUnusedImages returns a list of unused images
// Returns an error if the list cannot be retrieved
func (d *DockerClient) GetUnusedImages(olderThan int) ([]image.Summary, error) {
	// Get all images
	allImages, err := d.client.ImageList(d.ctx, image.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	// Get all containers
	containers, err := d.client.ContainerList(d.ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	// Map container image IDs
	usedImages := make(map[string]bool)
	for _, container := range containers {
		usedImages[container.ImageID] = true
	}

	// Filter unused images
	var unusedImages []image.Summary
	for _, image := range allImages {
		if !usedImages[image.ID] {
			// If olderThan is specified, check image age
			if olderThan > 0 {
				imageAge := time.Since(time.Unix(image.Created, 0))
				if imageAge.Hours() < float64(olderThan*24) {
					continue
				}
			}
			unusedImages = append(unusedImages, image)
		}
	}

	return unusedImages, nil
}

// PruneImages removes unused images
func (d *DockerClient) PruneImages(olderThan int) (image.PruneReport, error) {
	pruneFilters := filters.NewArgs()
	if olderThan > 0 {
		timestamp := time.Now().Add(-time.Hour * 24 * time.Duration(olderThan)).Format(time.RFC3339)
		pruneFilters.Add("until", timestamp)
	}
	return d.client.ImagesPrune(d.ctx, pruneFilters)
}

// GetDanglingImages gets the list of dangling images
func (d *DockerClient) GetDanglingImages() ([]image.Summary, error) {
	args := filters.NewArgs()
	args.Add("dangling", "true")

	return d.client.ImageList(d.ctx, image.ListOptions{Filters: args})
}

// PruneDanglingImages removes dangling images
func (d *DockerClient) PruneDanglingImages() (image.PruneReport, error) {
	pruneFilters := filters.NewArgs()
	pruneFilters.Add("dangling", "true")
	return d.client.ImagesPrune(d.ctx, pruneFilters)
}

// GetUnusedVolumes gets the list of unused volumes
func (d *DockerClient) GetUnusedVolumes() ([]volume.Volume, error) {
	// Get all volumes
	volumes, err := d.client.VolumeList(d.ctx, volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Get all containers to see which volumes are used
	containers, err := d.client.ContainerList(d.ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	// Map the names of used volumes
	usedVolumes := make(map[string]bool)
	for _, container := range containers {
		containerInfo, err := d.client.ContainerInspect(d.ctx, container.ID)
		if err != nil {
			continue
		}

		for _, mount := range containerInfo.Mounts {
			if mount.Type == "volume" {
				usedVolumes[mount.Name] = true
			}
		}
	}

	// Filter unused volumes
	var unusedVolumes []volume.Volume
	for _, volume := range volumes.Volumes {
		if !usedVolumes[volume.Name] {
			unusedVolumes = append(unusedVolumes, *volume)
		}
	}

	return unusedVolumes, nil
}

// PruneVolumes removes unused volumes
func (d *DockerClient) PruneVolumes() (volume.PruneReport, error) {
	pruneFilters := filters.NewArgs()
	return d.client.VolumesPrune(d.ctx, pruneFilters)
}

// GetUnusedNetworks gets the list of unused networks
func (d *DockerClient) GetUnusedNetworks() ([]network.Summary, error) {
	// Get all networks
	networks, err := d.client.NetworkList(d.ctx, network.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Get all containers to see which networks are used
	containers, err := d.client.ContainerList(d.ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	// Map the names of used networks
	usedNetworks := make(map[string]bool)
	for _, container := range containers {
		containerInfo, err := d.client.ContainerInspect(d.ctx, container.ID)
		if err != nil {
			continue
		}

		for networkName := range containerInfo.NetworkSettings.Networks {
			usedNetworks[networkName] = true
		}
	}

	// Add default networks that should not be removed
	defaultNetworks := []string{"bridge", "host", "none"}
	for _, network := range defaultNetworks {
		usedNetworks[network] = true
	}

	// Filter unused networks
	var unusedNetworks []network.Summary
	for _, network := range networks {
		if !usedNetworks[network.Name] {
			unusedNetworks = append(unusedNetworks, network)
		}
	}

	return unusedNetworks, nil
}

// PruneNetworks removes unused networks
func (d *DockerClient) PruneNetworks() (network.PruneReport, error) {
	pruneFilters := filters.NewArgs()
	return d.client.NetworksPrune(d.ctx, pruneFilters)
}

// GetUnusedBuilds gets the list of unused Docker builds
func (d *DockerClient) GetUnusedBuilds() ([]types.BuildCache, error) {
	// Get disk usage which includes build cache information
	diskUsage, err := d.GetDiskUsage()
	if err != nil {
		return nil, err
	}

	// Filter unused build caches
	var builds []types.BuildCache
	for _, cache := range diskUsage.BuildCache {
		if !cache.InUse {
			builds = append(builds, *cache)
		}
	}

	return builds, nil
}

// PruneBuilds removes unused Docker builds
func (d *DockerClient) PruneBuilds(olderThan int) (*types.BuildCachePruneReport, error) {
	pruneFilters := filters.NewArgs()

	if olderThan > 0 {
		timestamp := time.Now().Add(-time.Hour * 24 * time.Duration(olderThan)).Format(time.RFC3339)
		pruneFilters.Add("until", timestamp)
	}

	return d.client.BuildCachePrune(d.ctx, types.BuildCachePruneOptions{
		All:     true, // This will prune all unused build cache
		Filters: pruneFilters,
	})
}
