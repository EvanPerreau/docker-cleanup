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

// View represents the user interface
type View struct {
	YellowText func(format string, a ...interface{}) string
	GreenText  func(format string, a ...interface{}) string
	RedText    func(format string, a ...interface{}) string
}

// NewView creates a new View instance
func NewView() *View {
	return &View{
		YellowText: color.New(color.FgYellow).SprintfFunc(),
		GreenText:  color.New(color.FgGreen).SprintfFunc(),
		RedText:    color.New(color.FgRed).SprintfFunc(),
	}
}

// ShowTitle displays a colored title
func (v *View) ShowTitle(title string) {
	fmt.Println(v.YellowText(title))
}

// ShowSuccess displays a success message
func (v *View) ShowSuccess(message string) {
	fmt.Println(v.GreenText(message))
}

// ShowError displays an error message
func (v *View) ShowError(err error) {
	fmt.Println(v.RedText("Error: %v", err))
}

// ShowDiskUsage displays disk usage
func (v *View) ShowDiskUsage(diskUsage *types.DiskUsage) {
	v.ShowTitle("Current Docker disk usage:")

	fmt.Println("Containers :", len(diskUsage.Containers))
	fmt.Println("Images     :", len(diskUsage.Images))
	fmt.Println("Volumes    :", len(diskUsage.Volumes))
	fmt.Println("Builds     :", len(diskUsage.BuildCache))

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

	fmt.Printf("Total size: %s\n\n", FormatSize(totalSize))
}

// ShowContainers displays the list of containers
func (v *View) ShowContainers(containers []types.Container, dryRun bool) {
	if len(containers) == 0 {
		v.ShowSuccess("No stopped containers to remove.")
		return
	}

	fmt.Printf("Found %d containers to remove.\n", len(containers))

	if dryRun {
		v.ShowTitle("[DRY RUN] The following containers would be removed:")
		for _, container := range containers {
			fmt.Printf(" - %s (%s)\n", container.ID[:12], strings.Join(container.Names, ", "))
		}
	}
}

// ShowContainerRemoved displays a message for a removed container
func (v *View) ShowContainerRemoved(containerID string, names []string) {
	v.ShowSuccess(fmt.Sprintf("Container removed: %s (%s)", containerID[:12], strings.Join(names, ", ")))
}

// ShowContainersCleanupComplete displays a message for the end of container cleanup
func (v *View) ShowContainersCleanupComplete() {
	v.ShowSuccess("Stopped containers successfully removed.")
}

// ShowImages displays the list of images
func (v *View) ShowImages(images []image.Summary, dryRun bool, imageType string) {
	if len(images) == 0 {
		v.ShowSuccess(fmt.Sprintf("No %s images to remove.", imageType))
		return
	}

	if dryRun {
		v.ShowTitle(fmt.Sprintf("[DRY RUN] The following %s images would be removed:", imageType))
		for _, image := range images {
			tags := image.RepoTags
			if len(tags) == 0 {
				tags = []string{"<none>:<none>"}
			}
			fmt.Printf(" - %s (%s)\n", image.ID[:12], strings.Join(tags, ", "))
		}
	}
}

// ShowImagesPruneResult displays the result of image cleanup
func (v *View) ShowImagesPruneResult(report image.PruneReport, imageType string) {
	if len(report.ImagesDeleted) == 0 {
		v.ShowSuccess(fmt.Sprintf("No %s images to remove.", imageType))
		return
	}

	for _, img := range report.ImagesDeleted {
		if img.Untagged != "" {
			v.ShowSuccess(fmt.Sprintf("Image untagged: %s", img.Untagged))
		}
		if img.Deleted != "" {
			v.ShowSuccess(fmt.Sprintf("Image deleted: %s", img.Deleted))
		}
	}
	v.ShowSuccess(fmt.Sprintf("%s images successfully removed. Space reclaimed: %s", imageType, FormatSize(report.SpaceReclaimed)))
}

// ShowVolumes displays the list of volumes
func (v *View) ShowVolumes(volumes []volume.Volume, dryRun bool) {
	if len(volumes) == 0 {
		v.ShowSuccess("No unused volumes to remove.")
		return
	}

	if dryRun {
		v.ShowTitle("[DRY RUN] The following volumes would be removed:")
		for _, volume := range volumes {
			fmt.Printf(" - %s\n", volume.Name)
		}
	}
}

// ShowVolumesPruneResult displays the result of volume cleanup
func (v *View) ShowVolumesPruneResult(report volume.PruneReport) {
	if len(report.VolumesDeleted) == 0 {
		v.ShowSuccess("No unused volumes to remove.")
		return
	}

	for _, vol := range report.VolumesDeleted {
		v.ShowSuccess(fmt.Sprintf("Volume deleted: %s", vol))
	}
	v.ShowSuccess(fmt.Sprintf("Unused volumes successfully removed. Space reclaimed: %s", FormatSize(report.SpaceReclaimed)))
}

// ShowNetworks displays the list of networks
func (v *View) ShowNetworks(networks []network.Summary, dryRun bool) {
	if len(networks) == 0 {
		v.ShowSuccess("No unused networks to remove.")
		return
	}

	if dryRun {
		v.ShowTitle("[DRY RUN] The following networks would be removed:")
		for _, network := range networks {
			fmt.Printf(" - %s (%s)\n", network.Name, network.ID[:12])
		}
	}
}

// ShowNetworksPruneResult displays the result of network cleanup
func (v *View) ShowNetworksPruneResult(report network.PruneReport) {
	if len(report.NetworksDeleted) == 0 {
		v.ShowSuccess("No unused networks to remove.")
		return
	}

	for _, network := range report.NetworksDeleted {
		v.ShowSuccess(fmt.Sprintf("Network deleted: %s", network))
	}
	v.ShowSuccess("Unused networks successfully removed.")
}

// ShowCleanupComplete displays a message for the end of global cleanup
func (v *View) ShowCleanupComplete() {
	v.ShowSuccess("\nGlobal cleanup completed successfully!")
}

// FormatSize formats a size in bytes to a readable unit
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

// ShowBuilds displays the list of builds
func (v *View) ShowBuilds(builds []types.BuildCache, dryRun bool) {
	if len(builds) == 0 {
		fmt.Println("No builds to clean.")
		return
	}

	mode := "[DRY RUN] "
	if !dryRun {
		mode = ""
	}

	fmt.Printf("%sBuilds that would be removed (%d):\n", mode, len(builds))
	for _, build := range builds {
		fmt.Printf("- %s: %s\n", build.ID[:12], build.Description)
	}
}

// ShowBuildsPruneResult displays the result of builds cleanup
func (v *View) ShowBuildsPruneResult(report *types.BuildCachePruneReport) {
	if len(report.CachesDeleted) == 0 {
		fmt.Println("No builds were removed.")
		return
	}

	fmt.Printf("Builds removed (%d):\n", len(report.CachesDeleted))
	for _, id := range report.CachesDeleted {
		fmt.Printf("- %s\n", id[:12])
	}
	fmt.Printf("Space freed: %s\n", humanize.Bytes(uint64(report.SpaceReclaimed)))
}
