package controllers

import (
	"docker-cleanup/app/models"
	"docker-cleanup/app/views"
	"fmt"
)

// Config represents the configuration for the controller
type Config struct {
	DryRun    bool
	OlderThan int
	ShowSize  bool
}

// Controller manages interactions between the model and view
type Controller struct {
	model *models.DockerClient
	view  *views.View
}

// NewController creates a new Controller instance
func NewController() (*Controller, error) {
	dockerClient := models.NewDockerClient()

	return &Controller{
		model: dockerClient,
		view:  views.NewView(),
	}, nil
}

// Close closes the Docker client
func (c *Controller) Close() error {
	return c.model.Close()
}

// RunContainerCleanup executes the cleanup of containers
func (c *Controller) RunContainerCleanup() {
	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Removing stopped containers...")

	containers, err := c.model.GetStoppedContainers()
	if err != nil {
		c.view.ShowError(fmt.Errorf("error retrieving containers: %v", err))
		return
	}

	c.view.ShowContainers(containers, GetConfig().DryRun)

	if !GetConfig().DryRun && len(containers) > 0 {
		for _, container := range containers {
			if err := c.model.RemoveContainer(container.ID); err != nil {
				c.view.ShowError(fmt.Errorf("error removing container %s: %v", container.ID[:12], err))
			} else {
				c.view.ShowContainerRemoved(container.ID, container.Names)
			}
		}
		c.view.ShowContainersCleanupComplete()
	}
}

// RunImageCleanup executes the cleanup of unused images
func (c *Controller) RunImageCleanup() {
	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Removing unused images...")

	if GetConfig().DryRun {
		images, err := c.model.GetUnusedImages(GetConfig().OlderThan)
		if err != nil {
			c.view.ShowError(fmt.Errorf("error retrieving images: %v", err))
			return
		}

		c.view.ShowImages(images, true, "unused")
	} else {
		report, err := c.model.PruneImages(GetConfig().OlderThan)
		if err != nil {
			c.view.ShowError(fmt.Errorf("error removing images: %v", err))
			return
		}

		c.view.ShowImagesPruneResult(report, "unused")
	}
}

// RunDanglingCleanup executes the cleanup of dangling images
func (c *Controller) RunDanglingCleanup() {
	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Removing dangling images...")

	if GetConfig().DryRun {
		images, err := c.model.GetDanglingImages()
		if err != nil {
			c.view.ShowError(fmt.Errorf("error retrieving images: %v", err))
			return
		}

		c.view.ShowImages(images, true, "dangling")
	} else {
		report, err := c.model.PruneDanglingImages()
		if err != nil {
			c.view.ShowError(fmt.Errorf("error removing dangling images: %v", err))
			return
		}

		c.view.ShowImagesPruneResult(report, "dangling")
	}
}

// RunVolumeCleanup executes the cleanup of unused volumes
func (c *Controller) RunVolumeCleanup() {

	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Removing unused volumes...")

	if GetConfig().DryRun {
		volumes, err := c.model.GetUnusedVolumes()
		if err != nil {
			c.view.ShowError(fmt.Errorf("error retrieving volumes: %v", err))
			return
		}

		c.view.ShowVolumes(volumes, true)
	} else {
		report, err := c.model.PruneVolumes()
		if err != nil {
			c.view.ShowError(fmt.Errorf("error removing volumes: %v", err))
			return
		}

		c.view.ShowVolumesPruneResult(report)
	}
}

// RunNetworkCleanup executes the cleanup of unused networks
func (c *Controller) RunNetworkCleanup() {
	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Removing unused networks...")

	if GetConfig().DryRun {
		networks, err := c.model.GetUnusedNetworks()
		if err != nil {
			c.view.ShowError(fmt.Errorf("error retrieving networks: %v", err))
			return
		}

		c.view.ShowNetworks(networks, true)
	} else {
		report, err := c.model.PruneNetworks()
		if err != nil {
			c.view.ShowError(fmt.Errorf("error removing networks: %v", err))
			return
		}

		c.view.ShowNetworksPruneResult(report)
	}
}

// RunAllCleanup executes the cleanup of all Docker resources
func (c *Controller) RunAllCleanup() {
	c.RunContainerCleanup()
	fmt.Println()
	c.RunDanglingCleanup()
	fmt.Println()
	c.RunImageCleanup()
	fmt.Println()
	c.RunVolumeCleanup()
	fmt.Println()
	c.RunNetworkCleanup()
	fmt.Println()
	c.RunBuildsCleanup()

	c.view.ShowCleanupComplete()
}

// ShowDiskUsage displays current disk usage
func (c *Controller) ShowDiskUsage() {
	diskUsage, err := c.model.GetDiskUsage()
	if err != nil {
		c.view.ShowError(fmt.Errorf("failed to get disk usage: %v", err))
		return
	}

	c.view.ShowDiskUsage(diskUsage)
}

// RunBuildsCleanup executes the cleanup of unused Docker builds
func (c *Controller) RunBuildsCleanup() {
	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Removing Docker builds...")

	if GetConfig().DryRun {
		builds, err := c.model.GetUnusedBuilds()
		if err != nil {
			c.view.ShowError(fmt.Errorf("error retrieving builds: %v", err))
			return
		}

		c.view.ShowBuilds(builds, true)
	} else {
		report, err := c.model.PruneBuilds(GetConfig().OlderThan)
		if err != nil {
			c.view.ShowError(fmt.Errorf("error removing builds: %v", err))
			return
		}

		c.view.ShowBuildsPruneResult(report)
	}
}
