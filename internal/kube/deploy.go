package kube

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployCanary(clientset *kubernetes.Clientset) error {
	deploymentsClient := clientset.AppsV1().Deployments("default")

	// First check if it already exists
	existing, err := deploymentsClient.Get(context.TODO(), "nginx-canary", metav1.GetOptions{})
	if err == nil {
		fmt.Printf("ℹ️  Deployment '%s' already exists (created %v), skipping...\n", existing.Name, existing.CreationTimestamp)
		return nil
	}

	// If not found, continue creating it
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-canary",
			Labels: map[string]string{
				"app": "nginx-canary",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx-canary",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx-canary",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.25-alpine",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	fmt.Println("✅ Canary deployment 'nginx-canary' created.")
	return nil
}


func int32Ptr(i int32) *int32 {
	return &i
}
