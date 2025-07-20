package kube

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func StreamLogs(clientset *kubernetes.Clientset, deploymentName string) error {
	log.Printf("📜 Fetching logs for deployment '%s'...\n", deploymentName)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get pods in 'default' namespace with label selector matching the deployment
	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deploymentName),
	})
	if err != nil {
		return fmt.Errorf("failed to list pods: %w", err)
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("no pods found for deployment '%s'", deploymentName)
	}

	// Get logs from the first pod
	pod := pods.Items[0]
	req := clientset.CoreV1().Pods("default").GetLogs(pod.Name, &corev1.PodLogOptions{
		Follow: true,
	})

	stream, err := req.Stream(ctx)
	if err != nil {
		return fmt.Errorf("failed to stream logs: %w", err)
	}
	defer stream.Close()

	log.Printf("🪵 Streaming logs from pod: %s\n\n", pod.Name)
	_, err = io.Copy(log.Writer(), stream)
	return err
}
