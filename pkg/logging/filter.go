package logging

import (
	"bufio"
	"io"
)

func FilterAndFormatLogs(reader io.Reader, writer io.Writer, filterLevel LogLevel) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		entry := ParseLogEntry(scanner.Text())
		if entry.Level >= filterLevel {
			logLevelColors[entry.Level].Fprintf(writer, "%s\n", entry.Message)
		}
	}
	return scanner.Err()
}
