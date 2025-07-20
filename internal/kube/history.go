package kube

import (
	"context"
	"fmt"
	"sort"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func PrintDeploymentHistory(clientset *kubernetes.Clientset, name string) error {
	deploy, err := clientset.AppsV1().Deployments("default").Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	rsList, err := clientset.AppsV1().ReplicaSets("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list ReplicaSets: %w", err)
	}

	var history []appsv1.ReplicaSet
	for _, rs := range rsList.Items {
		if ownerRefs := rs.OwnerReferences; len(ownerRefs) > 0 && ownerRefs[0].UID == deploy.UID {
			history = append(history, rs)
		}
	}

	if len(history) == 0 {
		fmt.Println("ℹ️  No history found for deployment.")
		return nil
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].CreationTimestamp.Time.Before(history[j].CreationTimestamp.Time)
	})

	fmt.Printf("📜 Rollout history for deployment '%s':\n", name)
	fmt.Println("REVISION\tIMAGE\t\t\tCREATED")
	for i, rs := range history {
		image := "?"
		if len(rs.Spec.Template.Spec.Containers) > 0 {
			image = rs.Spec.Template.Spec.Containers[0].Image
		}
		timestamp := rs.CreationTimestamp.Local().Format(time.RFC1123)
		fmt.Printf("%d\t\t%s\t%s\n", i+1, image, timestamp)
	}

	return nil
}
