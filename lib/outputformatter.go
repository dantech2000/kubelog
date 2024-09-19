package lib

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// OutputFormatter handles different output formats
type OutputFormatter struct {
	PodName    string
	Namespace  string
	Containers []string
}

// NewOutputFormatter creates a new OutputFormatter
func NewOutputFormatter(podName, namespace string, containers []string) *OutputFormatter {
	return &OutputFormatter{
		PodName:    podName,
		Namespace:  namespace,
		Containers: containers,
	}
}

// FormatOutput formats the output based on the specified format
func (of *OutputFormatter) FormatOutput(format string) (string, error) {
	switch format {
	case "json":
		return of.formatJSON()
	case "yaml":
		return of.formatYAML()
	case "posix":
		return of.formatPOSIX()
	default:
		return FormatContainerList(of.PodName, of.Namespace, of.Containers), nil
	}
}

func (of *OutputFormatter) formatJSON() (string, error) {
	data := map[string]interface{}{
		"podName":    of.PodName,
		"namespace":  of.Namespace,
		"containers": of.Containers,
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling to JSON: %v", err)
	}
	return string(jsonData), nil
}

func (of *OutputFormatter) formatYAML() (string, error) {
	data := map[string]interface{}{
		"podName":    of.PodName,
		"namespace":  of.Namespace,
		"containers": of.Containers,
	}
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling to YAML: %v", err)
	}
	return string(yamlData), nil
}

func (of *OutputFormatter) formatPOSIX() (string, error) {
	return strings.Join(of.Containers, "\n"), nil
}
