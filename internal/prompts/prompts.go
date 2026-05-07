// Package prompts generates interactive prompts to fill in missing or placeholder
// values in a .env file, producing a completed map ready for export.
package prompts

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Options controls prompt behaviour.
type Options struct {
	// In is the reader used for user input (defaults to os.Stdin).
	In io.Reader
	// Out is the writer used to display prompts (defaults to os.Stdout).
	Out io.Writer
	// SkipFilled skips keys that already have a non-empty, non-placeholder value.
	SkipFilled bool
	// PlaceholderMarkers lists substrings that indicate a value is a placeholder.
	PlaceholderMarkers []string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		In:                 os.Stdin,
		Out:                os.Stdout,
		SkipFilled:         true,
		PlaceholderMarkers: []string{"<", "CHANGE_ME", "TODO", "YOUR_"},
	}
}

// Fill iterates over keys in the provided env map and prompts the user for
// values where needed. It returns a new map with the collected values merged
// on top of the originals.
func Fill(env map[string]string, opts Options) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = v
	}

	scanner := bufio.NewScanner(opts.In)

	for _, key := range sortedKeys(env) {
		val := env[key]
		if opts.SkipFilled && !needsInput(val, opts.PlaceholderMarkers) {
			continue
		}
		fmt.Fprintf(opts.Out, "Enter value for %s [%s]: ", key, val)
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			result[key] = input
		}
	}
	return result, scanner.Err()
}

func needsInput(val string, markers []string) bool {
	if val == "" {
		return true
	}
	for _, m := range markers {
		if strings.Contains(val, m) {
			return true
		}
	}
	return false
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// simple insertion sort for determinism without importing sort
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
