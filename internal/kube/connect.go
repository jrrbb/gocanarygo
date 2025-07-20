package kube

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func MustConnect() *kubernetes.Clientset {
	kubeconfig := flag.String("kubeconfig", defaultKubeconfigPath(), "Path to the kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("❌ Failed to build kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("❌ Failed to create Kubernetes client: %v", err)
	}

	log.Println("✅ Connected to Kubernetes cluster.")
	return clientset
}

func defaultKubeconfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("❌ Failed to find user home directory: %v", err)
	}
	return filepath.Join(home, ".kube", "config")
}
