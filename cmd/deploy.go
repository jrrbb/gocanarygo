package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/jrrbb/gocanarygo/internal/kube"
)

var (
	name     string
	image    string
	replicas int32
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

		err = kube.DeployCanary(clientset, name, image, replicas)
		if err != nil {
			fmt.Println("❌ Deployment failed:", err)
			return
		}

		fmt.Println("🚀 Deployment completed successfully.")
	},
}

func init() {
	deployCmd.Flags().StringVar(&name, "name", "nginx-canary", "Name of the deployment")
	deployCmd.Flags().StringVar(&image, "image", "nginx:1.25-alpine", "Container image to deploy")
	deployCmd.Flags().Int32Var(&replicas, "replicas", 1, "Number of pod replicas")

	rootCmd.AddCommand(deployCmd)
}
