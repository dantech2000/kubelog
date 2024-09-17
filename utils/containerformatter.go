package utils

import (
	"strings"

	"github.com/fatih/color"
)

// FormatContainerList formats the list of containers for pretty printing
func FormatContainerList(podName, namespace string, containers []string) string {
	var output strings.Builder

	// Print header
	headerColor := color.New(color.FgCyan, color.Bold)
	separatorColor := color.New(color.FgCyan)
	containerColor := color.New(color.FgYellow)

	headerColor.Fprintf(&output, "\nContainers in pod '%s' (namespace: %s):\n", podName, namespace)
	separatorColor.Fprintln(&output, strings.Repeat("=", 50))

	// Print containers
	for _, container := range containers {
		containerColor.Fprintln(&output, container)
	}

	// Print footer
	separatorColor.Fprintln(&output, strings.Repeat("=", 50))
	headerColor.Fprintf(&output, "Total containers: %d\n\n", len(containers))

	return output.String()
}
