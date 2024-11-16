package logging

import (
	"testing"
	"time"
)

func TestParseLogEntry(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected LogEntry
	}{
		{
			name:  "Plain text log with level",
			input: "2024-03-15 12:19:57 DEBUG Starting application",
			expected: LogEntry{
				Level:   DEBUG,
				Message: "2024-03-15 12:19:57 DEBUG Starting application",
				Format:  FormatPlainText,
			},
		},
		{
			name:  "Plain text log with level",
			input: "2024-03-15 12:19:57 INFO Starting application",
			expected: LogEntry{
				Level:   INFO,
				Message: "2024-03-15 12:19:57 INFO Starting application",
				Format:  FormatPlainText,
			},
		},
		{
			name:  "Plain text log with level",
			input: "2024-03-15 12:19:57 WARN Starting application",
			expected: LogEntry{
				Level:   WARN,
				Message: "2024-03-15 12:19:57 WARN Starting application",
				Format:  FormatPlainText,
			},
		},
		{
			name:  "Plain text log with level",
			input: "2024-03-15 12:19:57 ERROR Starting application",
			expected: LogEntry{
				Level:   ERROR,
				Message: "2024-03-15 12:19:57 ERROR Starting application",
				Format:  FormatPlainText,
			},
		},
		{
			name:  "JSON log with Logrus format",
			input: `{"level":"info","msg":"Server started","time":"2024-03-15T12:19:57Z","port":8080}`,
			expected: LogEntry{
				Level:   INFO,
				Message: "Server started",
				Format:  FormatJSON,
				Logger:  "logrus",
				Fields: map[string]interface{}{
					"level": "info",
					"msg":   "Server started",
					"time":  "2024-03-15T12:19:57Z",
					"port":  float64(8080),
				},
			},
		},
		{
			name:  "JSON log with Zap format",
			input: `{"level":"error","ts":1647340797,"caller":"api/handler.go:42","msg":"Failed to process request","error":"invalid input"}`,
			expected: LogEntry{
				Level:   ERROR,
				Message: "Failed to process request",
				Format:  FormatJSON,
				Logger:  "zap",
				Fields: map[string]interface{}{
					"level":  "error",
					"ts":     float64(1647340797),
					"caller": "api/handler.go:42",
					"msg":    "Failed to process request",
					"error":  "invalid input",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLogEntry(tt.input)

			// Compare basic fields
			if got.Level != tt.expected.Level {
				t.Errorf("Level = %v, want %v", got.Level, tt.expected.Level)
			}
			if got.Format != tt.expected.Format {
				t.Errorf("Format = %v, want %v", got.Format, tt.expected.Format)
			}
			if got.Logger != tt.expected.Logger {
				t.Errorf("Logger = %v, want %v", got.Logger, tt.expected.Logger)
			}

			// For JSON logs, compare fields
			if tt.expected.Format == FormatJSON {
				if len(got.Fields) != len(tt.expected.Fields) {
					t.Errorf("Fields count = %v, want %v", len(got.Fields), len(tt.expected.Fields))
				}
				for k, v := range tt.expected.Fields {
					if got.Fields[k] != v {
						t.Errorf("Fields[%q] = %v, want %v", k, got.Fields[k], v)
					}
				}
			}
		})
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected LogLevel
		wantErr  bool
	}{
		{"debug level", "DEBUG", DEBUG, false},
		{"info level", "INFO", INFO, false},
		{"warn level", "WARN", WARN, false},
		{"warning level", "WARNING", WARN, false},
		{"error level", "ERROR", ERROR, false},
		{"lowercase level", "debug", DEBUG, false},
		{"mixed case level", "InFo", INFO, false},
		{"numeric debug level", "10", DEBUG, false},
		{"numeric info level", "20", INFO, false},
		{"numeric warn level", "30", WARN, false},
		{"numeric error level", "40", ERROR, false},
		{"invalid level", "INVALID", DEBUG, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLogLevel(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLogLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("ParseLogLevel() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDetectLogger(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			name: "Zap logger",
			input: map[string]interface{}{
				"caller": "main.go:42",
				"ts":     1647340797,
			},
			expected: "zap",
		},
		{
			name: "Bunyan logger",
			input: map[string]interface{}{
				"@level":     "info",
				"@timestamp": "2024-03-15T12:19:57Z",
			},
			expected: "bunyan",
		},
		{
			name: "Logrus logger",
			input: map[string]interface{}{
				"level": "info",
				"msg":   "test message",
			},
			expected: "logrus",
		},
		{
			name: "Unknown logger",
			input: map[string]interface{}{
				"custom": "field",
			},
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectLogger(tt.input)
			if got != tt.expected {
				t.Errorf("detectLogger() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantHour  int
		wantMin   int
		wantSec   int
	}{
		{
			name:      "RFC3339",
			input:     "2024-03-15T12:19:57Z",
			wantError: false,
			wantYear:  2024,
			wantMonth: 3,
			wantDay:   15,
			wantHour:  12,
			wantMin:   19,
			wantSec:   57,
		},
		{
			name:      "RFC3339 with timezone",
			input:     "2024-03-15T12:19:57+00:00",
			wantError: false,
			wantYear:  2024,
			wantMonth: 3,
			wantDay:   15,
			wantHour:  12,
			wantMin:   19,
			wantSec:   57,
		},
		{
			name:      "RFC3339 with nanoseconds",
			input:     "2024-03-15T12:19:57.123456789Z",
			wantError: false,
			wantYear:  2024,
			wantMonth: 3,
			wantDay:   15,
			wantHour:  12,
			wantMin:   19,
			wantSec:   57,
		},
		{
			name:      "Simple date time",
			input:     "2024-03-15 12:19:57",
			wantError: false,
			wantYear:  2024,
			wantMonth: 3,
			wantDay:   15,
			wantHour:  12,
			wantMin:   19,
			wantSec:   57,
		},
		{
			name:      "Date with slashes",
			input:     "2024/03/15 12:19:57",
			wantError: false,
			wantYear:  2024,
			wantMonth: 3,
			wantDay:   15,
			wantHour:  12,
			wantMin:   19,
			wantSec:   57,
		},
		{
			name:      "Invalid timestamp",
			input:     "not a timestamp",
			wantError: true,
		},
		{
			name:      "Empty timestamp",
			input:     "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTimestamp(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("parseTimestamp() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if got.Year() != tt.wantYear {
					t.Errorf("Year = %v, want %v", got.Year(), tt.wantYear)
				}
				if got.Month() != tt.wantMonth {
					t.Errorf("Month = %v, want %v", got.Month(), tt.wantMonth)
				}
				if got.Day() != tt.wantDay {
					t.Errorf("Day = %v, want %v", got.Day(), tt.wantDay)
				}
				if got.Hour() != tt.wantHour {
					t.Errorf("Hour = %v, want %v", got.Hour(), tt.wantHour)
				}
				if got.Minute() != tt.wantMin {
					t.Errorf("Minute = %v, want %v", got.Minute(), tt.wantMin)
				}
				if got.Second() != tt.wantSec {
					t.Errorf("Second = %v, want %v", got.Second(), tt.wantSec)
				}
			}
		})
	}
}
