package cmd

import (
	"docker-cleanup/app/controllers"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var buildsCmd = &cobra.Command{
	Use:   "builds",
	Short: "Clean unused Docker builds",
	Long:  `Removes Docker builds that are no longer used in the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctrl, err := controllers.NewController()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer ctrl.Close()

		ctrl.RunBuildsCleanup()
	},
}
