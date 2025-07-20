package cmd

import (
	"log"

	"gocanarygo/internal/kube"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs <deployment-name>",
	Short: "Stream logs from the first pod in a deployment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		clientset := kube.MustConnect()

		if err := kube.StreamLogs(clientset, name); err != nil {
			log.Fatalf("❌ Error streaming logs: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
