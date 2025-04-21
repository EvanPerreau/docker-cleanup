package cmd

import (
	"docker-cleanup/app/controllers"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Nettoie les ressources Docker inutilisées",
	Long:  `Un outil CLI pour nettoyer facilement les conteneurs arrêtés, les images inutilisées, les volumes et les réseaux dans Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctrl, err := controllers.NewController()
		if err != nil {
			fmt.Printf("Erreur: %v\n", err)
			os.Exit(1)
		}
		defer ctrl.Close()

		ctrl.RunAllCleanup()
	},
}
