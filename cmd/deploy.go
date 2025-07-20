package cmd

import (
	"log"

	"gocanarygo/internal/kube"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a Kubernetes workload",
	Run: func(cmd *cobra.Command, args []string) {
		clientset := kube.MustConnect()

		name, _ := cmd.Flags().GetString("name")
		image, _ := cmd.Flags().GetString("image")
		replicas, _ := cmd.Flags().GetInt32("replicas")

		err := kube.CreateDeployment(clientset, name, image, replicas)
		if err != nil {
			log.Fatalf("❌ Deployment failed: %v", err)
		}
	},
}

func init() {
	deployCmd.Flags().String("name", "nginx", "Name of the deployment")
	deployCmd.Flags().String("image", "nginx", "Container image")
	deployCmd.Flags().Int32("replicas", 2, "Number of replicas")

	rootCmd.AddCommand(deployCmd)
}
