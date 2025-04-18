package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "docker-cleanup",
	Short: "Docker Cleanup Tool - Nettoie les ressources Docker inutilisées",
	Long:  `Un outil CLI pour nettoyer facilement les conteneurs arrêtés, les images inutilisées, les volumes et les réseaux dans Docker.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
