package cmd

import (
	"log"

	"gocanarygo/internal/kube"
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback --name <deployment-name>",
	Short: "Rollback a deployment to its previous revision",
	Run: func(cmd *cobra.Command, args []string) {
		clientset := kube.MustConnect()

		if err := kube.RollbackDeployment(clientset, rollbackName); err != nil {
			log.Fatalf("❌ Rollback failed: %v", err)
		}
	},
}

var rollbackName string

func init() {
	rollbackCmd.Flags().StringVar(&rollbackName, "name", "", "Deployment name to rollback")
	rollbackCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(rollbackCmd)
}
