package cmd

import (
	"log"

	"gocanarygo/internal/kube"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history <deployment-name>",
	Short: "Show rollout history of a deployment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		clientset := kube.MustConnect()

		if err := kube.PrintDeploymentHistory(clientset, name); err != nil {
			log.Fatalf("❌ Failed to show history: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
