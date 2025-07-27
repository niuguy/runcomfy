package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of runcomfy",
	Long:  `Print the version number and build information for runcomfy.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("runcomfy v1.0.0")
		fmt.Println("ComfyUI workflow analyzer and dependency manager")
		fmt.Println("Built for RunPod environments")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}