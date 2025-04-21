package cmd

import (
	"docker-cleanup/app/controllers"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "docker-cleanup",
	Short: "Docker Cleanup Tool - Clean unused Docker resources",
	Long:  `A CLI tool to easily clean stopped containers, unused images, volumes, and networks in Docker.`,
}

func Execute() {
	rootCmd.PersistentFlags().BoolVar(&controllers.GetConfig().DryRun, "dry-run", false, "Run in dry run mode (default: false)")
	rootCmd.PersistentFlags().IntVar(&controllers.GetConfig().OlderThan, "older-than", 0, "Keep resources older than N days (default: 0)")
	rootCmd.PersistentFlags().BoolVar(&controllers.GetConfig().ShowSize, "show-size", false, "Show size of resources (default: false)")

	rootCmd.AddCommand(containersCmd)
	rootCmd.AddCommand(imagesCmd)
	rootCmd.AddCommand(networksCmd)
	rootCmd.AddCommand(volumesCmd)
	rootCmd.AddCommand(danglingImagesCmd)
	rootCmd.AddCommand(allCmd)
	rootCmd.AddCommand(buildsCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
