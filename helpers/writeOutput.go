package helpers

import (
	"bytes"
	"log/slog"
	"os"
)

// WriteOutput turns bytes into a file
func WriteOutput(filename string, b bytes.Buffer) (int64, error) {
	f, err := os.Create(filename)
	if err != nil {
		slog.Error("Error creating file:", "error", err)
	}
	defer f.Close()
	n, err := b.WriteTo(f)
	if err != nil {
		slog.Error("error writing to file:", "error", err)
	}
	return n, err
}
