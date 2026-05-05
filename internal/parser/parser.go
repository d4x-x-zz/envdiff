package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Env represents a parsed .env file as a map of key-value pairs.
type Env map[string]string

// ParseFile reads a .env file from the given path and returns an Env map.
// Lines starting with '#' are treated as comments and skipped.
// Empty lines are also skipped.
func ParseFile(path string) (Env, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer f.Close()

	env := make(Env)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a single KEY=VALUE line into its components.
func parseLine(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid line %q: expected KEY=VALUE format", line)
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return "", "", fmt.Errorf("empty key in line %q", line)
	}

	value := strings.TrimSpace(parts[1])
	// Strip surrounding quotes if present
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}

	return key, value, nil
}
