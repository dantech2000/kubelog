package logging

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

type LogEntry struct {
	Level   LogLevel
	Message string
}

var logLevelColors = map[LogLevel]*color.Color{
	DEBUG: color.New(color.FgCyan),
	INFO:  color.New(color.FgGreen),
	WARN:  color.New(color.FgYellow),
	ERROR: color.New(color.FgRed),
}

func ParseLogEntry(line string) LogEntry {
	levelRegex := regexp.MustCompile(`(?i)\b(DEBUG|INFO|WARN(?:ING)?|ERROR)\b`)
	match := levelRegex.FindString(line)

	level := DEBUG // Default to DEBUG if no level is found
	if match != "" {
		parsedLevel, _ := ParseLogLevel(match)
		level = parsedLevel
	}

	return LogEntry{
		Level:   level,
		Message: line,
	}
}

func ParseLogLevel(level string) (LogLevel, error) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	default:
		return DEBUG, fmt.Errorf("invalid log level: %s", level)
	}
}

func ParseLog(log string) string {
	entry := ParseLogEntry(log)
	return logLevelColors[entry.Level].Sprint(log)
}
