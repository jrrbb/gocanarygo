package cmd

import (
	"log"

	"gocanarygo/internal/kube"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all deployments in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		clientset := kube.MustConnect()
		if err := kube.ListDeployments(clientset); err != nil {
			log.Fatalf("❌ Error listing deployments: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
