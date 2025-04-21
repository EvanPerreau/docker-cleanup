package cmd

import (
	"docker-cleanup/app/controllers"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Clean unused images",
	Long:  `Cleans unused images in Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctrl, err := controllers.NewController()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer ctrl.Close()

		ctrl.RunImageCleanup()
	},
}
