package cmd

import (
	"fmt"

	"gocanarygo/internal/kube"
	"github.com/spf13/cobra"
)

var (
	hpaName       string
	hpaMin        int32
	hpaMax        int32
	hpaCPUPercent int32
)

var autoscaleCmd = &cobra.Command{
	Use:   "autoscale",
	Short: "Set up Horizontal Pod Autoscaler for a deployment",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := kube.GetClientSet()
		if err != nil {
			fmt.Println("❌ Failed to connect to cluster:", err)
			return
		}

		err = kube.SetupHPA(clientset, hpaName, hpaMin, hpaMax, hpaCPUPercent)
		if err != nil {
			fmt.Println("❌ Failed to set up HPA:", err)
			return
		}

		fmt.Println("🚀 Autoscaler ready.")
	},
}

func init() {
	autoscaleCmd.Flags().StringVar(&hpaName, "name", "", "Deployment name")
	autoscaleCmd.Flags().Int32Var(&hpaMin, "min", 1, "Minimum replicas")
	autoscaleCmd.Flags().Int32Var(&hpaMax, "max", 5, "Maximum replicas")
	autoscaleCmd.Flags().Int32Var(&hpaCPUPercent, "cpu", 80, "Target CPU utilization (%)")

	autoscaleCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(autoscaleCmd)
}
