package kube

import (
	"context"
	"fmt"
	"log"
	"sort"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RollbackDeployment(clientset *kubernetes.Clientset, name string) error {
	log.Printf("🔙 Rolling back deployment '%s'...\n", name)

	deployments := clientset.AppsV1().Deployments("default")
	replicaSets := clientset.AppsV1().ReplicaSets("default")

	// Get current deployment
	deploy, err := deployments.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	// List all ReplicaSets owned by this deployment
	rsList, err := replicaSets.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list ReplicaSets: %w", err)
	}

	var ownedRS []appsv1.ReplicaSet
	for _, rs := range rsList.Items {
		for _, owner := range rs.OwnerReferences {
			if owner.Kind == "Deployment" && owner.Name == deploy.Name {
				ownedRS = append(ownedRS, rs)
			}
		}
	}

	// Sort by creation timestamp descending
	sort.Slice(ownedRS, func(i, j int) bool {
		return ownedRS[i].CreationTimestamp.After(ownedRS[j].CreationTimestamp.Time)
	})

	if len(ownedRS) < 2 {
		return fmt.Errorf("not enough ReplicaSets to perform rollback")
	}

	previousRS := ownedRS[1]

	// Rollback by updating the deployment's template to match the previous RS
	deploy.Spec.Template = previousRS.Spec.Template

	_, err = deployments.Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	log.Printf("✅ Rolled back deployment '%s' to previous revision\n", name)
	return nil
}
