package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/jrrbb/gocanarygo/internal/kube"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a canary workload to your Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := kube.GetClientSet()
		if err != nil {
			fmt.Println("❌ Failed to connect to cluster:", err)
			return
		}

		err = kube.DeployCanary(clientset)
		if err != nil {
			fmt.Println("❌ Deployment failed:", err)
			return
		}

		fmt.Println("🚀 Deployment completed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
