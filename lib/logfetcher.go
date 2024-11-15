package lib

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type LogFetcher struct {
	Clientset     *kubernetes.Clientset
	Namespace     string
	PodName       string
	ContainerName string
	Follow        bool
	Previous      bool
	Writer        io.Writer
}

func NewLogFetcher(clientset *kubernetes.Clientset, namespace, podName string, follow bool, previous bool, writer io.Writer) *LogFetcher {
	return &LogFetcher{
		Clientset: clientset,
		Namespace: namespace,
		PodName:   podName,
		Follow:    follow,
		Previous:  previous,
		Writer:    writer,
	}
}

func (lf *LogFetcher) GetLogs() error {
	// Get container name first if not specified
	if lf.ContainerName == "" {
		containerName, err := lf.getSingleContainerName()
		if err != nil {
			return err
		}
		lf.ContainerName = containerName
	}

	// Now proceed with log fetching
	podLogOpts := corev1.PodLogOptions{
		Container: lf.ContainerName,
		Follow:    lf.Follow,
		Previous:  lf.Previous,
	}

	req := lf.Clientset.CoreV1().Pods(lf.Namespace).GetLogs(lf.PodName, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return fmt.Errorf("error opening log stream: %v", err)
	}
	defer podLogs.Close()

	_, err = io.Copy(lf.Writer, podLogs)
	if err != nil {
		return fmt.Errorf("error copying log stream: %v", err)
	}
	fmt.Printf("Log stream opened successfully\n")
	return nil
}

type containerInfo struct {
	Name       string
	Ready      bool
	Status     string
	Image      string
	DisplayStr string
}

func getContainerStatus(pod *corev1.Pod, containerName string) (bool, string) {
	for _, status := range pod.Status.ContainerStatuses {
		if status.Name == containerName {
			return status.Ready, getContainerState(status.State)
		}
	}
	return false, "Unknown"
}

func getContainerState(state corev1.ContainerState) string {
	if state.Running != nil {
		return "Running"
	}
	if state.Waiting != nil {
		return fmt.Sprintf("Waiting (%s)", state.Waiting.Reason)
	}
	if state.Terminated != nil {
		return fmt.Sprintf("Terminated (%s)", state.Terminated.Reason)
	}
	return "Unknown"
}

func formatContainerInfo(info containerInfo) string {
	statusColor := color.New(color.FgRed)
	if info.Ready {
		statusColor = color.New(color.FgGreen)
	}

	readySymbol := "✗"
	if info.Ready {
		readySymbol = "✓"
	}

	return fmt.Sprintf("%s %s [%s] (%s)",
		statusColor.Sprint(readySymbol),
		info.Name,
		info.Status,
		info.Image)
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
	}

	// Create container info list for the prompt
	containers := make([]containerInfo, containerCount)
	options := make([]string, containerCount)

	for i, c := range pod.Spec.Containers {
		ready, status := getContainerStatus(pod, c.Name)
		info := containerInfo{
			Name:   c.Name,
			Ready:  ready,
			Status: status,
			Image:  c.Image,
		}
		containers[i] = info
		options[i] = formatContainerInfo(info)
	}

	// Prepare the survey prompt
	var selectedIdx int
	prompt := &survey.Select{
		Message: "Choose a container:",
		Options: options,
		Filter: func(filter string, value string, index int) bool {
			container := containers[index]
			filter = strings.ToLower(filter)
			return strings.Contains(strings.ToLower(container.Name), filter) ||
				strings.Contains(strings.ToLower(container.Status), filter) ||
				strings.Contains(strings.ToLower(container.Image), filter)
		},
	}

	// Show the prompt and get user's selection
	err = survey.AskOne(prompt, &selectedIdx, survey.WithPageSize(10))
	if err != nil {
		if err == terminal.InterruptErr {
			return "", fmt.Errorf("operation cancelled")
		}
		return "", fmt.Errorf("selection failed: %v", err)
	}

	return containers[selectedIdx].Name, nil
}

func GetKubernetesClient() (*kubernetes.Clientset, string, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, "", err
	}

	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		return nil, "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, "", err
	}

	return clientset, namespace, nil
}

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
