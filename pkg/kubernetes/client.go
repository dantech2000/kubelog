// Package kubernetes provides functionality for interacting with Kubernetes clusters
package kubernetes

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubernetesClient creates a new Kubernetes client using the default kubeconfig.
// It returns the clientset, the current namespace, and any error encountered.
// The current namespace is determined from the kubeconfig context.
func GetKubernetesClient() (*kubernetes.Clientset, string, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get namespace from config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return clientset, namespace, nil
}
