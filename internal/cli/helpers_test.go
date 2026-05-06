package cli

import (
	"bytes"
	"io"
	"os"
)

// captureOutput redirects stdout during fn execution and returns the captured string.
func captureOutput(fn func() error) (string, error) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	fnErr := fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return "", err
	}
	return buf.String(), fnErr
}
