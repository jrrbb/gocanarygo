package kube

import (
	"context"
	"fmt"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func PrintDeploymentStatus(clientset *kubernetes.Clientset, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get Deployment
	deploy, err := clientset.AppsV1().Deployments("default").Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("could not get deployment: %w", err)
	}

	printDeployment(deploy)

	// Get HPA (optional)
	hpa, err := clientset.AutoscalingV2().HorizontalPodAutoscalers("default").Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Printf("ℹ️  No autoscaler found for '%s' (or failed to fetch it)", name)
		return nil
	}

	printHPA(hpa)
	return nil
}

func printDeployment(deploy *appsv1.Deployment) {
	log.Println("📦 Deployment Info:")
	log.Printf("  Name:        %s\n", deploy.Name)
	log.Printf("  Replicas:    %d desired / %d available\n", *deploy.Spec.Replicas, deploy.Status.AvailableReplicas)
	log.Printf("  Image:       %s\n", deploy.Spec.Template.Spec.Containers[0].Image)
	log.Printf("  Created At:  %s\n", deploy.CreationTimestamp.Time.Format(time.RFC1123))
}

func printHPA(hpa *autoscalingv2.HorizontalPodAutoscaler) {
	log.Println("📊 Autoscaler Info:")
	if hpa.Spec.MinReplicas != nil {
		log.Printf("  Min Pods:    %d\n", *hpa.Spec.MinReplicas)
	}
	log.Printf("  Max Pods:    %d\n", hpa.Spec.MaxReplicas)

	// Print CPU metric info
	for i, metric := range hpa.Spec.Metrics {
		if metric.Type == autoscalingv2.ResourceMetricSourceType && metric.Resource != nil {
			target := "unknown"
			if metric.Resource.Target.AverageUtilization != nil {
				target = fmt.Sprintf("%d%%", *metric.Resource.Target.AverageUtilization)
			}
			log.Printf("  Target CPU Utilization (spec metric #%d): %s", i+1, target)
		}
	}

	for i, metric := range hpa.Status.CurrentMetrics {
		if metric.Type == autoscalingv2.ResourceMetricSourceType && metric.Resource != nil {
			current := "unknown"
			if metric.Resource.Current.AverageUtilization != nil {
				current = fmt.Sprintf("%d%%", *metric.Resource.Current.AverageUtilization)
			}
			log.Printf("  Current CPU Utilization (status metric #%d): %s", i+1, current)
		}
	}
}
