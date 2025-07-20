package kube

import (
	"context"
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreateDeployment(clientset *kubernetes.Clientset, name, image string, replicas int32) error {
	deploymentsClient := clientset.AppsV1().Deployments(corev1.NamespaceDefault)

	labels := map[string]string{"app": name}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: 80},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(80),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       10,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(80),
									},
								},
								InitialDelaySeconds: 2,
								PeriodSeconds:       5,
							},
						},
					},
				},
			},
		},
	}

	log.Printf("📦 Creating deployment '%s' with image '%s' and %d replicas...\n", name, image, replicas)
	_, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	// Create Horizontal Pod Autoscaler
	hpaClient := clientset.AutoscalingV2().HorizontalPodAutoscalers(corev1.NamespaceDefault)
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       name,
			},
			MinReplicas: &replicas,
			MaxReplicas: 5,
			Metrics: []autoscalingv2.MetricSpec{
				{
					Type: autoscalingv2.ResourceMetricSourceType,
					Resource: &autoscalingv2.ResourceMetricSource{
						Name: corev1.ResourceCPU,
						Target: autoscalingv2.MetricTarget{
							Type:               autoscalingv2.UtilizationMetricType,
							AverageUtilization: int32Ptr(50),
						},
					},
				},
			},
		},
	}

	log.Printf("📈 Creating HPA for deployment '%s'...\n", name)
	_, err = hpaClient.Create(context.TODO(), hpa, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create HPA: %w", err)
	}

	log.Println("🚀 Deployment completed successfully.")
	return nil
}

func int32Ptr(i int32) *int32 {
	return &i
}
