package lib

// LogLevel represents the severity of a log entry
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// LogEntry represents a single log entry
type LogEntry struct {
	Level   LogLevel
	Message string
}

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR"}[l]
}
