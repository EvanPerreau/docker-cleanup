package cmd

import (
	"docker-cleanup/app/controllers"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var volumesCmd = &cobra.Command{
	Use:   "volumes",
	Short: "Clean unused volumes",
	Long:  `Cleans unused volumes in Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctrl, err := controllers.NewController()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer ctrl.Close()

		ctrl.RunVolumeCleanup()
	},
}
