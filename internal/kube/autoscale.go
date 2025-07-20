package kube

import (
	"context"
	"fmt"
	"log"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func SetupHPA(clientset *kubernetes.Clientset, name string, minReplicas, maxReplicas int32, targetCPU int32) error {
	hpaClient := clientset.AutoscalingV2().HorizontalPodAutoscalers("default")

	// Define CPU metric
	cpuMetric := autoscalingv2.MetricSpec{
		Type: autoscalingv2.ResourceMetricSourceType,
		Resource: &autoscalingv2.ResourceMetricSource{
			Name: "cpu",
			Target: autoscalingv2.MetricTarget{
				Type:               autoscalingv2.UtilizationMetricType,
				AverageUtilization: &targetCPU,
			},
		},
	}

	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       name,
			},
			MinReplicas: &minReplicas,
			MaxReplicas: maxReplicas,
			Metrics:     []autoscalingv2.MetricSpec{cpuMetric},
		},
	}

	// Try to update if it exists
	existing, err := hpaClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err == nil {
		hpa.ResourceVersion = existing.ResourceVersion
		_, err = hpaClient.Update(context.TODO(), hpa, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update HPA: %w", err)
		}
		log.Printf("🔄 Updated HPA for '%s'", name)
		return nil
	}

	// Create if it doesn't exist
	_, err = hpaClient.Create(context.TODO(), hpa, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create HPA: %w", err)
	}

	log.Printf("✅ HPA created for '%s' (min: %d, max: %d, cpu: %d%%)", name, minReplicas, maxReplicas, targetCPU)
	return nil
}
