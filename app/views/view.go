package views

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

// View représente l'interface utilisateur
type View struct {
	YellowText func(format string, a ...interface{}) string
	GreenText  func(format string, a ...interface{}) string
	RedText    func(format string, a ...interface{}) string
}

// NewView crée une nouvelle instance de View
func NewView() *View {
	return &View{
		YellowText: color.New(color.FgYellow).SprintfFunc(),
		GreenText:  color.New(color.FgGreen).SprintfFunc(),
		RedText:    color.New(color.FgRed).SprintfFunc(),
	}
}

// ShowTitle affiche un titre coloré
func (v *View) ShowTitle(title string) {
	fmt.Println(v.YellowText(title))
}

// ShowSuccess affiche un message de succès
func (v *View) ShowSuccess(message string) {
	fmt.Println(v.GreenText(message))
}

// ShowError affiche un message d'erreur
func (v *View) ShowError(err error) {
	fmt.Println(v.RedText("Erreur: %v", err))
}

// ShowDiskUsage affiche l'utilisation du disque
func (v *View) ShowDiskUsage(diskUsage *types.DiskUsage) {
	v.ShowTitle("Taille actuelle utilisée par Docker :")

	fmt.Println("Conteneurs :", len(diskUsage.Containers))
	fmt.Println("Images      :", len(diskUsage.Images))
	fmt.Println("Volumes     :", len(diskUsage.Volumes))
	fmt.Println("Builds      :", len(diskUsage.BuildCache))

	var totalSize uint64
	for _, container := range diskUsage.Containers {
		totalSize += uint64(container.SizeRw)
	}
	for _, image := range diskUsage.Images {
		totalSize += uint64(image.Size)
	}
	for _, volume := range diskUsage.Volumes {
		totalSize += uint64(volume.UsageData.Size)
	}
	for _, build := range diskUsage.BuildCache {
		totalSize += uint64(build.Size)
	}

	fmt.Printf("Taille totale : %s\n\n", FormatSize(totalSize))
}

// ShowContainers affiche la liste des conteneurs
func (v *View) ShowContainers(containers []types.Container, dryRun bool) {
	if len(containers) == 0 {
		v.ShowSuccess("Aucun conteneur arrêté à supprimer.")
		return
	}

	fmt.Printf("Trouvé %d conteneurs à supprimer.\n", len(containers))

	if dryRun {
		v.ShowTitle("[DRY RUN] Les conteneurs suivants seraient supprimés:")
		for _, container := range containers {
			fmt.Printf(" - %s (%s)\n", container.ID[:12], strings.Join(container.Names, ", "))
		}
	}
}

// ShowContainerRemoved affiche un message pour un conteneur supprimé
func (v *View) ShowContainerRemoved(containerID string, names []string) {
	v.ShowSuccess(fmt.Sprintf("Conteneur supprimé: %s (%s)", containerID[:12], strings.Join(names, ", ")))
}

// ShowContainersCleanupComplete affiche un message pour la fin du nettoyage des conteneurs
func (v *View) ShowContainersCleanupComplete() {
	v.ShowSuccess("Conteneurs arrêtés supprimés avec succès.")
}

// ShowImages affiche la liste des images
func (v *View) ShowImages(images []image.Summary, dryRun bool, imageType string) {
	if len(images) == 0 {
		v.ShowSuccess(fmt.Sprintf("Aucune image %s à supprimer.", imageType))
		return
	}

	if dryRun {
		v.ShowTitle(fmt.Sprintf("[DRY RUN] Les images %s suivantes seraient supprimées:", imageType))
		for _, image := range images {
			tags := image.RepoTags
			if len(tags) == 0 {
				tags = []string{"<none>:<none>"}
			}
			fmt.Printf(" - %s (%s)\n", image.ID[:12], strings.Join(tags, ", "))
		}
	}
}

// ShowImagesPruneResult affiche le résultat du nettoyage des images
func (v *View) ShowImagesPruneResult(report image.PruneReport, imageType string) {
	if len(report.ImagesDeleted) == 0 {
		v.ShowSuccess(fmt.Sprintf("Aucune image %s à supprimer.", imageType))
		return
	}

	for _, img := range report.ImagesDeleted {
		if img.Untagged != "" {
			v.ShowSuccess(fmt.Sprintf("Image untagged: %s", img.Untagged))
		}
		if img.Deleted != "" {
			v.ShowSuccess(fmt.Sprintf("Image supprimée: %s", img.Deleted))
		}
	}
	v.ShowSuccess(fmt.Sprintf("Images %s supprimées avec succès. Espace récupéré: %s", imageType, FormatSize(report.SpaceReclaimed)))
}

// ShowVolumes affiche la liste des volumes
func (v *View) ShowVolumes(volumes []volume.Volume, dryRun bool) {
	if len(volumes) == 0 {
		v.ShowSuccess("Aucun volume non utilisé à supprimer.")
		return
	}

	if dryRun {
		v.ShowTitle("[DRY RUN] Les volumes suivants seraient supprimés:")
		for _, volume := range volumes {
			fmt.Printf(" - %s\n", volume.Name)
		}
	}
}

// ShowVolumesPruneResult affiche le résultat du nettoyage des volumes
func (v *View) ShowVolumesPruneResult(report volume.PruneReport) {
	if len(report.VolumesDeleted) == 0 {
		v.ShowSuccess("Aucun volume non utilisé à supprimer.")
		return
	}

	for _, vol := range report.VolumesDeleted {
		v.ShowSuccess(fmt.Sprintf("Volume supprimé: %s", vol))
	}
	v.ShowSuccess(fmt.Sprintf("Volumes non utilisés supprimés avec succès. Espace récupéré: %s", FormatSize(report.SpaceReclaimed)))
}

// ShowNetworks affiche la liste des réseaux
func (v *View) ShowNetworks(networks []network.Summary, dryRun bool) {
	if len(networks) == 0 {
		v.ShowSuccess("Aucun réseau non utilisé à supprimer.")
		return
	}

	if dryRun {
		v.ShowTitle("[DRY RUN] Les réseaux suivants seraient supprimés:")
		for _, network := range networks {
			fmt.Printf(" - %s (%s)\n", network.Name, network.ID[:12])
		}
	}
}

// ShowNetworksPruneResult affiche le résultat du nettoyage des réseaux
func (v *View) ShowNetworksPruneResult(report network.PruneReport) {
	if len(report.NetworksDeleted) == 0 {
		v.ShowSuccess("Aucun réseau non utilisé à supprimer.")
		return
	}

	for _, network := range report.NetworksDeleted {
		v.ShowSuccess(fmt.Sprintf("Réseau supprimé: %s", network))
	}
	v.ShowSuccess("Réseaux non utilisés supprimés avec succès.")
}

// ShowCleanupComplete affiche un message pour la fin du nettoyage global
func (v *View) ShowCleanupComplete() {
	v.ShowSuccess("\nNettoyage global terminé avec succès !")
}

// FormatSize formate une taille en octets en une unité lisible
func FormatSize(size uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// ShowBuilds affiche la liste des builds
func (v *View) ShowBuilds(builds []types.BuildCache, dryRun bool) {
	if len(builds) == 0 {
		fmt.Println("Aucun build à nettoyer.")
		return
	}

	mode := "[DRY RUN] "
	if !dryRun {
		mode = ""
	}

	fmt.Printf("%sBuilds qui seraient supprimés (%d):\n", mode, len(builds))
	for _, build := range builds {
		fmt.Printf("- %s: %s\n", build.ID[:12], build.Description)
	}
}

// ShowBuildsPruneResult affiche le résultat du nettoyage des builds
func (v *View) ShowBuildsPruneResult(report *types.BuildCachePruneReport) {
	if len(report.CachesDeleted) == 0 {
		fmt.Println("Aucun build n'a été supprimé.")
		return
	}

	fmt.Printf("Builds supprimés (%d):\n", len(report.CachesDeleted))
	for _, id := range report.CachesDeleted {
		fmt.Printf("- %s\n", id[:12])
	}
	fmt.Printf("Espace libéré: %s\n", humanize.Bytes(uint64(report.SpaceReclaimed)))
}
