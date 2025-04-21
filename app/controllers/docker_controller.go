package controllers

import (
	"docker-cleanup/app/models"
	"docker-cleanup/app/views"
	"fmt"
)

// Config représente la configuration pour le contrôleur
type Config struct {
	DryRun    bool
	OlderThan int
	ShowSize  bool
}

// Controller gère les interactions entre le modèle et la vue
type Controller struct {
	model *models.DockerClient
	view  *views.View
}

// NewController crée une nouvelle instance de Controller
func NewController() (*Controller, error) {
	dockerClient := models.NewDockerClient()

	return &Controller{
		model: dockerClient,
		view:  views.NewView(),
	}, nil
}

// Close ferme le client Docker
func (c *Controller) Close() error {
	return c.model.Close()
}

// RunContainerCleanup exécute le nettoyage des conteneurs
func (c *Controller) RunContainerCleanup() {
	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Suppression des conteneurs arrêtés...")

	containers, err := c.model.GetStoppedContainers()
	if err != nil {
		c.view.ShowError(fmt.Errorf("erreur lors de la récupération des conteneurs: %v", err))
		return
	}

	c.view.ShowContainers(containers, GetConfig().DryRun)

	if !GetConfig().DryRun && len(containers) > 0 {
		for _, container := range containers {
			if err := c.model.RemoveContainer(container.ID); err != nil {
				c.view.ShowError(fmt.Errorf("erreur lors de la suppression du conteneur %s: %v", container.ID[:12], err))
			} else {
				c.view.ShowContainerRemoved(container.ID, container.Names)
			}
		}
		c.view.ShowContainersCleanupComplete()
	}
}

// RunImageCleanup exécute le nettoyage des images non utilisées
func (c *Controller) RunImageCleanup() {
	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Suppression des images non utilisées...")

	if GetConfig().DryRun {
		images, err := c.model.GetUnusedImages(GetConfig().OlderThan)
		if err != nil {
			c.view.ShowError(fmt.Errorf("erreur lors de la récupération des images: %v", err))
			return
		}

		c.view.ShowImages(images, true, "non utilisées")
	} else {
		report, err := c.model.PruneImages(GetConfig().OlderThan)
		if err != nil {
			c.view.ShowError(fmt.Errorf("erreur lors de la suppression des images: %v", err))
			return
		}

		c.view.ShowImagesPruneResult(report, "non utilisées")
	}
}

// RunDanglingCleanup exécute le nettoyage des images dangling
func (c *Controller) RunDanglingCleanup() {
	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Suppression des images dangling...")

	if GetConfig().DryRun {
		images, err := c.model.GetDanglingImages()
		if err != nil {
			c.view.ShowError(fmt.Errorf("erreur lors de la récupération des images: %v", err))
			return
		}

		c.view.ShowImages(images, true, "dangling")
	} else {
		report, err := c.model.PruneDanglingImages()
		if err != nil {
			c.view.ShowError(fmt.Errorf("erreur lors de la suppression des images dangling: %v", err))
			return
		}

		c.view.ShowImagesPruneResult(report, "dangling")
	}
}

// RunVolumeCleanup exécute le nettoyage des volumes non utilisés
func (c *Controller) RunVolumeCleanup() {

	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Suppression des volumes non utilisés...")

	if GetConfig().DryRun {
		volumes, err := c.model.GetUnusedVolumes()
		if err != nil {
			c.view.ShowError(fmt.Errorf("erreur lors de la récupération des volumes: %v", err))
			return
		}

		c.view.ShowVolumes(volumes, true)
	} else {
		report, err := c.model.PruneVolumes()
		if err != nil {
			c.view.ShowError(fmt.Errorf("erreur lors de la suppression des volumes: %v", err))
			return
		}

		c.view.ShowVolumesPruneResult(report)
	}
}

// RunNetworkCleanup exécute le nettoyage des réseaux non utilisés
func (c *Controller) RunNetworkCleanup() {
	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Suppression des réseaux non utilisés...")

	if GetConfig().DryRun {
		networks, err := c.model.GetUnusedNetworks()
		if err != nil {
			c.view.ShowError(fmt.Errorf("erreur lors de la récupération des réseaux: %v", err))
			return
		}

		c.view.ShowNetworks(networks, true)
	} else {
		report, err := c.model.PruneNetworks()
		if err != nil {
			c.view.ShowError(fmt.Errorf("erreur lors de la suppression des réseaux: %v", err))
			return
		}

		c.view.ShowNetworksPruneResult(report)
	}
}

// RunAllCleanup exécute le nettoyage de toutes les ressources Docker
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

// ShowDiskUsage affiche l'utilisation du disque actuelle
func (c *Controller) ShowDiskUsage() {
	diskUsage, err := c.model.GetDiskUsage()
	if err != nil {
		c.view.ShowError(fmt.Errorf("échec de l'obtention de l'utilisation du disque: %v", err))
		return
	}

	c.view.ShowDiskUsage(diskUsage)
}

// RunBuildsCleanup exécute le nettoyage des builds Docker non utilisés
func (c *Controller) RunBuildsCleanup() {
	if GetConfig().ShowSize {
		c.ShowDiskUsage()
	}

	c.view.ShowTitle("Suppression des builds Docker...")

	if GetConfig().DryRun {
		builds, err := c.model.GetUnusedBuilds()
		if err != nil {
			c.view.ShowError(fmt.Errorf("erreur lors de la récupération des builds: %v", err))
			return
		}

		c.view.ShowBuilds(builds, true)
	} else {
		report, err := c.model.PruneBuilds(GetConfig().OlderThan)
		if err != nil {
			c.view.ShowError(fmt.Errorf("erreur lors de la suppression des builds: %v", err))
			return
		}

		c.view.ShowBuildsPruneResult(report)
	}
}
