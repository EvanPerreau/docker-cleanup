package cmd

import (
	"docker-cleanup/app/controllers"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var buildsCmd = &cobra.Command{
	Use:   "builds",
	Short: "Nettoie les builds Docker non utilisés",
	Long:  `Supprime les builds Docker qui ne sont plus utilisés dans le système.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctrl, err := controllers.NewController()
		if err != nil {
			fmt.Printf("Erreur: %v\n", err)
			os.Exit(1)
		}
		defer ctrl.Close()

		ctrl.RunBuildsCleanup()
	},
}
