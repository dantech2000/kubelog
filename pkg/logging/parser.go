package logging

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// LogFormat represents different log formats we can handle
type LogFormat int

const (
	FormatPlainText LogFormat = iota
	FormatJSON
)

// LogEntry represents a parsed log entry with all possible fields
type LogEntry struct {
	Level     LogLevel
	Message   string
	Timestamp time.Time
	Fields    map[string]interface{}
	Format    LogFormat
	Logger    string // Logger type (e.g., "zap", "logrus", etc.)
}

var logLevelColors = map[LogLevel]*color.Color{
	DEBUG: color.New(color.FgCyan),
	INFO:  color.New(color.FgGreen),
	WARN:  color.New(color.FgYellow),
	ERROR: color.New(color.FgRed),
}

// Colors for different log components
var (
	timestampColor = color.New(color.FgBlue)
	loggerColor    = color.New(color.FgMagenta)
	keyColor       = color.New(color.FgCyan)
	valueColor     = color.New(color.FgWhite)
	quoteColor     = color.New(color.FgHiBlack)
	errorColor     = color.New(color.FgRed, color.Bold)
)

// Common field mappings for different JSON log formats
var (
	// Level field names across different loggers
	jsonLevelFields = []string{
		"level",     // Common
		"severity",  // Google Cloud
		"log_level", // Custom
		"loglevel",  // Custom
		"@level",    // Bunyan
		"levelname", // Python logging
		"status",    // NGINX
		"LEVEL",     // Some uppercase variants
	}

	// Message field names
	jsonMessageFields = []string{
		"message",  // Common
		"msg",      // Zap, Logrus
		"log",      // Docker
		"text",     // Custom
		"@message", // Bunyan
		"MESSAGE",  // Systemd
	}

	// Timestamp field names
	jsonTimeFields = []string{
		"time",       // Common
		"timestamp",  // Common
		"@timestamp", // ELK
		"ts",         // Zap
		"Time",       // AWS CloudWatch
		"TIME",       // Some uppercase variants
		"datetime",   // Python logging
	}

	// Time formats to try parsing
	timeFormats = []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05.000Z0700",
		"2006-01-02T15:04:05Z0700",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		time.UnixDate,
		time.ANSIC,
	}
)

// detectLogFormat determines if the log line is JSON or plain text
func detectLogFormat(line string) LogFormat {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "{") && strings.HasSuffix(line, "}") {
		var js map[string]interface{}
		if json.Unmarshal([]byte(line), &js) == nil {
			return FormatJSON
		}
	}
	return FormatPlainText
}

// detectLogger tries to determine which logging framework generated the log
func detectLogger(data map[string]interface{}) string {
	switch {
	case data["caller"] != nil && data["ts"] != nil:
		return "zap"
	case data["@level"] != nil && data["@timestamp"] != nil:
		return "bunyan"
	case data["log.level"] != nil:
		return "winston"
	case data["levelname"] != nil:
		return "python"
	case data["stream"] != nil && data["log"] != nil:
		return "docker"
	case data["level"] != nil && data["msg"] != nil:
		return "logrus"
	default:
		return "unknown"
	}
}

// parseTimestamp attempts to parse a timestamp string using various formats
func parseTimestamp(timeStr string) (time.Time, error) {
	for _, format := range timeFormats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

// parseJSONLog attempts to parse a JSON log entry
func parseJSONLog(line string) LogEntry {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(line), &data); err != nil {
		return LogEntry{
			Level:   DEBUG,
			Message: line,
			Format:  FormatJSON,
		}
	}

	logger := detectLogger(data)
	entry := LogEntry{
		Format: FormatJSON,
		Fields: data,
		Logger: logger,
	}

	// Find and parse level
	for _, field := range jsonLevelFields {
		if val, ok := data[field]; ok {
			// Handle both string and numeric levels
			levelStr := fmt.Sprintf("%v", val)
			if level, err := ParseLogLevel(levelStr); err == nil {
				entry.Level = level
				break
			}
		}
	}

	// Find message
	for _, field := range jsonMessageFields {
		if val, ok := data[field]; ok {
			entry.Message = fmt.Sprintf("%v", val)
			break
		}
	}

	// If no message found, try error field or full line
	if entry.Message == "" {
		if err, ok := data["error"]; ok {
			entry.Message = fmt.Sprintf("%v", err)
		} else {
			entry.Message = line
		}
	}

	// Parse timestamp
	for _, field := range jsonTimeFields {
		if val, ok := data[field]; ok {
			// Handle numeric timestamps (milliseconds since epoch)
			if numTime, ok := val.(float64); ok {
				msec := int64(numTime)
				if msec > 1e11 { // Assuming milliseconds if number is large enough
					entry.Timestamp = time.Unix(0, msec*int64(time.Millisecond))
				} else {
					entry.Timestamp = time.Unix(msec, 0)
				}
				break
			}
			// Try parsing string timestamps
			if timeStr, ok := val.(string); ok {
				if ts, err := parseTimestamp(timeStr); err == nil {
					entry.Timestamp = ts
					break
				}
			}
		}
	}

	return entry
}

// parsePlainTextLog parses a plain text log entry
func parsePlainTextLog(line string) LogEntry {
	levelRegex := regexp.MustCompile(`(?i)\b(DEBUG|INFO|WARN(?:ING)?|ERROR|FATAL|TRACE)\b`)
	match := levelRegex.FindString(line)

	level := DEBUG // Default to DEBUG if no level is found
	if match != "" {
		parsedLevel, _ := ParseLogLevel(match)
		level = parsedLevel
	}

	// Try to find and parse timestamp at the beginning of the line
	timeRegex := regexp.MustCompile(`^\d{4}[-/]\d{2}[-/]\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})?`)
	if timeStr := timeRegex.FindString(line); timeStr != "" {
		if ts, err := parseTimestamp(timeStr); err == nil {
			return LogEntry{
				Level:     level,
				Message:   line,
				Format:    FormatPlainText,
				Timestamp: ts,
			}
		}
	}

	return LogEntry{
		Level:   level,
		Message: line,
		Format:  FormatPlainText,
	}
}

func ParseLogEntry(line string) LogEntry {
	format := detectLogFormat(line)
	if format == FormatJSON {
		return parseJSONLog(line)
	}
	return parsePlainTextLog(line)
}

// ParseLogLevel parses both string and numeric log levels
func ParseLogLevel(level string) (LogLevel, error) {
	// Try parsing numeric levels first
	if numLevel, err := strconv.Atoi(level); err == nil {
		switch {
		case numLevel <= 10:
			return DEBUG, nil // DEBUG
		case numLevel <= 20:
			return INFO, nil // INFO
		case numLevel <= 30:
			return WARN, nil // WARN
		case numLevel > 30:
			return ERROR, nil // ERROR
		}
	}

	normalizedLevel := strings.ToUpper(level)
	// Handle common variations
	switch {
	case strings.HasPrefix(normalizedLevel, "DEBUG") || normalizedLevel == "TRACE" || normalizedLevel == "FINE":
		return DEBUG, nil
	case strings.HasPrefix(normalizedLevel, "INFO") || normalizedLevel == "NOTICE":
		return INFO, nil
	case strings.HasPrefix(normalizedLevel, "WARN"):
		return WARN, nil
	case strings.HasPrefix(normalizedLevel, "ERR") || normalizedLevel == "CRITICAL" || normalizedLevel == "FATAL":
		return ERROR, nil
	default:
		return DEBUG, fmt.Errorf("invalid log level: %s", level)
	}
}

func formatLogEntry(entry LogEntry) string {
	var parts []string

	// Add timestamp if available
	if !entry.Timestamp.IsZero() {
		parts = append(parts, timestampColor.Sprintf("[%s]", entry.Timestamp.Format("2006-01-02 15:04:05")))
	}

	// Add level with appropriate color
	parts = append(parts, logLevelColors[entry.Level].Sprintf("[%s]", entry.Level))

	// Add logger type for JSON logs with a distinctive style
	if entry.Format == FormatJSON && entry.Logger != "unknown" {
		parts = append(parts, loggerColor.Sprintf("[%s]", entry.Logger))
	}

	// Format message based on content
	message := entry.Message
	if entry.Level == ERROR {
		message = errorColor.Sprint(message)
	} else if strings.Contains(message, "error") || strings.Contains(message, "failed") {
		// Highlight error-related messages even if not marked as ERROR level
		message = errorColor.Sprint(message)
	}
	parts = append(parts, message)

	// For JSON logs, format fields with enhanced visibility
	if entry.Format == FormatJSON {
		excludeFields := map[string]bool{
			"level": true, "severity": true, "log_level": true,
			"message": true, "msg": true, "text": true,
			"time": true, "timestamp": true, "@timestamp": true,
			"logger": true, "caller": true,
		}

		var fields []string
		for k, v := range entry.Fields {
			if !excludeFields[k] {
				formattedValue := formatValue(v)
				fields = append(fields, fmt.Sprintf("%s=%s",
					keyColor.Sprint(k),
					formattedValue))
			}
		}

		if len(fields) > 0 {
			parts = append(parts, strings.Join(fields, " "))
		}
	}

	return strings.Join(parts, " ")
}

// formatValue formats a value with appropriate coloring based on its type
func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		if val == "" {
			return quoteColor.Sprint(`""`)
		}
		// Check if the string needs quotes (contains spaces or special characters)
		if strings.ContainsAny(val, " =,\"'[]{}()") {
			return fmt.Sprintf("%s%s%s",
				quoteColor.Sprint(`"`),
				valueColor.Sprint(val),
				quoteColor.Sprint(`"`))
		}
		return valueColor.Sprint(val)
	case nil:
		return quoteColor.Sprint("null")
	case bool:
		if val {
			return valueColor.Sprint("true")
		}
		return valueColor.Sprint("false")
	case float64:
		if float64(int64(val)) == val {
			return valueColor.Sprintf("%d", int64(val))
		}
		return valueColor.Sprintf("%.2f", val)
	default:
		return valueColor.Sprintf("%v", val)
	}
}

func ParseLog(log string) string {
	entry := ParseLogEntry(log)
	return formatLogEntry(entry)
}
