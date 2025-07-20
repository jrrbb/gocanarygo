package kube

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CleanupDeployment(clientset *kubernetes.Clientset, name string) error {
	ctx := context.TODO()

	// Delete HPA (ignore errors if not found)
	if err := clientset.AutoscalingV1().HorizontalPodAutoscalers("default").Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		log.Printf("⚠️  Failed to delete HPA (might not exist): %v", err)
	} else {
		log.Printf("🧨 HPA '%s' deleted.", name)
	}

	// Delete deployment
	err := clientset.AppsV1().Deployments("default").Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete deployment: %v", err)
	}

	log.Printf("🧹 Deployment '%s' deleted.", name)
	return nil
}
