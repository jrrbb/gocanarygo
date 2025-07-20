package cmd

import (
	"log"

	"gocanarygo/internal/kube"
	"github.com/spf13/cobra"
)

var (
	scaleName     string
	scaleReplicas int32
)

var scaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale an existing deployment to a specific number of replicas",
	Run: func(cmd *cobra.Command, args []string) {
		clientset := kube.MustConnect()

		if err := kube.ScaleDeployment(clientset, scaleName, scaleReplicas); err != nil {
			log.Fatalf("❌ Error scaling deployment: %v", err)
		}
	},
}

func init() {
	scaleCmd.Flags().StringVar(&scaleName, "name", "", "Deployment name to scale")
	scaleCmd.Flags().Int32Var(&scaleReplicas, "replicas", 1, "Desired number of replicas")

	scaleCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(scaleCmd)
}
