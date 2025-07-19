package kube

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DeployCanary(clientset *kubernetes.Clientset, name, image string, replicas int32) error {
	deploymentsClient := clientset.AppsV1().Deployments("default")

	// Try to get the existing deployment
	existing, err := deploymentsClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err == nil {
		// Deployment exists – update it
		existing.Spec.Replicas = int32Ptr(replicas)
		existing.Spec.Template.Spec.Containers[0].Image = image

		_, err = deploymentsClient.Update(context.TODO(), existing, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update deployment '%s': %w", name, err)
		}

		fmt.Printf("🔄 Deployment '%s' updated to image '%s' with %d replicas.\n", name, image, replicas)
		return nil
	}

	// Deployment does not exist – create it
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: 80},
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

	fmt.Printf("✅ Deployment '%s' with image '%s' created with %d replicas.\n", name, image, replicas)
	return nil
}


func int32Ptr(i int32) *int32 {
	return &i
}
