package cmd

import (
	"gocanarygo/internal/kube"
	"github.com/spf13/cobra"
	"log"
)

var describeCmd = &cobra.Command{
	Use:   "describe [deployment-name]",
	Short: "Describe a deployment in detail",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deploymentName := args[0]
		log.Println("✅ Connected to Kubernetes cluster.")
		err := kube.DescribeDeployment(deploymentName)
		if err != nil {
			log.Fatalf("❌ Describe failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
