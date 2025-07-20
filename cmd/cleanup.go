package cmd

import (
	"log"
	"os"

	"gocanarygo/internal/kube"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup --name [deployment-name]",
	Short: "Delete a Kubernetes deployment and any attached autoscaler",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			log.Println("❌ Please provide a deployment name using --name")
			os.Exit(1)
		}
		clientset := kube.MustConnect()
		err := kube.CleanupDeployment(clientset, name)
		if err != nil {
			log.Println("❌ Cleanup error:", err)
			os.Exit(1)
		}
		log.Println("✅ Cleanup complete.")
	},
}

func init() {
	cleanupCmd.Flags().String("name", "", "Deployment name")
	rootCmd.AddCommand(cleanupCmd)
}
