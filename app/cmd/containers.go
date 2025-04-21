package cmd

import (
	"docker-cleanup/app/controllers"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var containersCmd = &cobra.Command{
	Use:   "containers",
	Short: "Clean stopped containers",
	Long:  `Cleans stopped containers in Docker.`,
	Run: func(cmd *cobra.Command, args []string) {

		ctrl, err := controllers.NewController()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer ctrl.Close()

		ctrl.RunContainerCleanup()
	},
}
