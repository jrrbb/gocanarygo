package kube

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ScaleDeployment scales the specified deployment to the desired number of replicas.
func ScaleDeployment(clientset *kubernetes.Clientset, name string, replicas int32) error {
	log.Printf("📏 Scaling deployment '%s' to %d replicas...\n", name, replicas)

	deploymentsClient := clientset.AppsV1().Deployments("default")

	// Fetch the current deployment
	deployment, err := deploymentsClient.Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	// Update the replicas
	deployment.Spec.Replicas = &replicas

	_, err = deploymentsClient.Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to scale deployment: %w", err)
	}

	log.Printf("✅ Deployment '%s' successfully scaled to %d replicas.\n", name, replicas)
	return nil
}
