package kube

import (
	"context"
	"fmt"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListDeployments(clientset *kubernetes.Clientset) error {
	log.Println("📋 Listing deployments in 'default' namespace...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deploymentsClient := clientset.AppsV1().Deployments("default")
	list, err := deploymentsClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	if len(list.Items) == 0 {
		log.Println("ℹ️  No deployments found.")
		return nil
	}

	for _, d := range list.Items {
		log.Printf("• %s (%d replicas) — image: %s\n",
			d.Name,
			*getReplicas(d),
			getImage(d),
		)
	}

	return nil
}

func getReplicas(d appsv1.Deployment) *int32 {
	if d.Spec.Replicas != nil {
		return d.Spec.Replicas
	}
	var zero int32 = 0
	return &zero
}

func getImage(d appsv1.Deployment) string {
	if len(d.Spec.Template.Spec.Containers) > 0 {
		return d.Spec.Template.Spec.Containers[0].Image
	}
	return "unknown"
}
