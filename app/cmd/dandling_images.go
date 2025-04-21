package cmd

import (
	"docker-cleanup/app/controllers"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var danglingImagesCmd = &cobra.Command{
	Use:   "dangling-images",
	Short: "Nettoie les images inutilisées",
	Long:  `Nettoie les images inutilisées dans Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctrl, err := controllers.NewController()
		if err != nil {
			fmt.Printf("Erreur: %v\n", err)
			os.Exit(1)
		}
		defer ctrl.Close()

		ctrl.RunDanglingCleanup()
	},
}
