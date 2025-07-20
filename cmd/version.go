package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "v0.1.0" // change this when you tag releases

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of gocanarygo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gocanarygo version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
