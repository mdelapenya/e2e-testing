package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eod",
	Short: "eod (Elastic Observability Dev) makes it easier to develop Observability projects.",
	Long: `A Fast and Flexible CLI for developing and testing Elastic's Observability projects
				built with love by mdelapenya and friends in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
