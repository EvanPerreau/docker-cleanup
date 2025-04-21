package cmd

import (
	"docker-cleanup/app/controllers"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "docker-cleanup",
	Short: "Docker Cleanup Tool - Nettoie les ressources Docker inutilisées",
	Long:  `Un outil CLI pour nettoyer facilement les conteneurs arrêtés, les images inutilisées, les volumes et les réseaux dans Docker.`,
}

func Execute() {
	rootCmd.PersistentFlags().BoolVar(&controllers.GetConfig().DryRun, "dry-run", false, "Run in dry run mode (default: false)")
	rootCmd.PersistentFlags().IntVar(&controllers.GetConfig().OlderThan, "older-than", 0, "Keep resources older than N days (default: 0)")
	rootCmd.PersistentFlags().BoolVar(&controllers.GetConfig().ShowSize, "show-size", false, "Show size of resources (default: false)")

	rootCmd.AddCommand(containersCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
