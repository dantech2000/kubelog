package lib

import (
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type LogFetcher struct {
	Clientset *kubernetes.Clientset
	Namespace string
	PodName   string
	Container string
	Follow    bool
}

func NewLogFetcher(clientset *kubernetes.Clientset, namespace, podName, container string, follow bool) *LogFetcher {
	return &LogFetcher{
		Clientset: clientset,
		Namespace: namespace,
		PodName:   podName,
		Container: container,
		Follow:    follow,
	}
}

func (lf *LogFetcher) FetchLogs() error {
	fmt.Printf("Fetching logs for pod %s in namespace %s, container %s\n", lf.PodName, lf.Namespace, lf.Container)

	// If container name is not provided, check if pod has only one container
	if lf.Container == "" {
		container, err := lf.getSingleContainerName()
		if err != nil {
			return err
		}
		lf.Container = container
	}

	// Fetch logs
	logOptions := &corev1.PodLogOptions{
		Follow:    lf.Follow,
		Container: lf.Container,
	}

	req := lf.Clientset.CoreV1().Pods(lf.Namespace).GetLogs(lf.PodName, logOptions)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return fmt.Errorf("error opening log stream: %v", err)
	}
	defer podLogs.Close()

	buffer := make([]byte, 2000)
	for {
		numBytes, err := podLogs.Read(buffer)
		if numBytes > 0 {
			logLine := string(buffer[:numBytes])
			coloredLog := ParseLog(logLine)
			fmt.Print(coloredLog)
		}
		if err != nil {
			break
		}
	}
	fmt.Printf("Log stream opened successfully\n")
	return nil
}

func (lf *LogFetcher) getSingleContainerName() (string, error) {
	pod, err := lf.Clientset.CoreV1().Pods(lf.Namespace).Get(context.Background(), lf.PodName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("error fetching pod details: %v", err)
	}

	containerCount := len(pod.Spec.Containers)
	if containerCount == 0 {
		return "", fmt.Errorf("no containers found in pod %s", lf.PodName)
	} else if containerCount == 1 {
		return pod.Spec.Containers[0].Name, nil
	} else {
		containerNames := make([]string, containerCount)
		for i, c := range pod.Spec.Containers {
			containerNames[i] = c.Name
		}
		return "", fmt.Errorf("pod %s has multiple containers. Please specify one using the --container flag: %v", lf.PodName, containerNames)
	}
}

// GetKubernetesClient initializes and returns a Kubernetes clientset
func GetKubernetesClient() (*kubernetes.Clientset, error) {
	kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// Add this function to the existing file

func ListContainers(clientset *kubernetes.Clientset, namespace, podName string) ([]string, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error fetching pod details: %v", err)
	}

	containers := make([]string, len(pod.Spec.Containers))
	for i, container := range pod.Spec.Containers {
		containers[i] = container.Name
	}

	return containers, nil
}

func (lf *LogFetcher) GetLogReader() (io.ReadCloser, error) {
	return lf.Clientset.CoreV1().Pods(lf.Namespace).GetLogs(lf.PodName, &corev1.PodLogOptions{
		Container: lf.Container,
		Follow:    lf.Follow,
	}).Stream(context.Background())
}
