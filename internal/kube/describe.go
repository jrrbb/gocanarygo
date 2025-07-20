package kube

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getKubernetesClient() (*kubernetes.Clientset, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
			return nil, fmt.Errorf("kubeconfig not found at %s", kubeconfig)
		}
	} else {
		return nil, fmt.Errorf("cannot determine home directory")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func DescribeDeployment(name string) error {
	clientset, err := getKubernetesClient()
	if err != nil {
		return err
	}

	deploy, err := clientset.AppsV1().Deployments("default").Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	printDeploymentDescribe(deploy)
	return nil
}

func printDeploymentDescribe(d *appsv1.Deployment) {
	log.Println("🧾 Deployment Description:")
	log.Printf("  Name:       %s", d.Name)
	log.Printf("  Namespace:  %s", d.Namespace)
	log.Printf("  Labels:     %s", mapToString(d.Labels))
	log.Printf("  Replicas:   %d desired / %d updated / %d available", *d.Spec.Replicas, d.Status.UpdatedReplicas, d.Status.AvailableReplicas)
	log.Printf("  Created:    %s", d.CreationTimestamp.Local().Format(time.RFC1123))

	log.Println("📦 Containers:")
	for _, c := range d.Spec.Template.Spec.Containers {
		log.Printf("  - Name:     %s", c.Name)
		log.Printf("    Image:    %s", c.Image)
		log.Printf("    Ports:    %v", c.Ports)
		log.Printf("    Env:      %v", c.Env)
	}
}

func mapToString(m map[string]string) string {
	if len(m) == 0 {
		return "none"
	}
	var b strings.Builder
	for k, v := range m {
		b.WriteString(fmt.Sprintf("%s=%s ", k, v))
	}
	return strings.TrimSpace(b.String())
}
