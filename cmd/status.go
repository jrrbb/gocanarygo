package cmd

import (
	"log"
	"os"

	"gocanarygo/internal/kube"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [deployment-name]",
	Short: "Check the status of a Kubernetes deployment and its autoscaler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deployName := args[0]
		clientset := kube.MustConnect()
		err := kube.PrintDeploymentStatus(clientset, deployName)
		if err != nil {
			log.Println("❌ Error checking status:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
