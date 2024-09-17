package utils

import (
	"regexp"

	"github.com/fatih/color"
)

var logLevels = map[string]*color.Color{
	"ERROR": color.New(color.FgRed),
	"WARN":  color.New(color.FgYellow),
	"INFO":  color.New(color.FgGreen),
}

func ParseLog(log string) string {
	for level, colorFunc := range logLevels {
		re := regexp.MustCompile(`(?i)\b` + level + `\b`)
		if re.MatchString(log) {
			return colorFunc.Sprint(log)
		}
	}
	return log
}

func AddLogLevel(level string, c *color.Color) {
	logLevels[level] = c
}
