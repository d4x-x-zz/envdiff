package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Format represents the output format for exported env templates.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatMarkdown Format = "markdown"
)

// Options controls how the template is exported.
type Options struct {
	Format      Format
	IncludeVals bool // if false, values are redacted
}

// Export writes an env template derived from the diff result to the given file path.
// Keys present in either side are included; values are redacted unless opts.IncludeVals is true.
func Export(result *differ.Result, leftEnv map[string]string, path string, opts Options) error {
	keys := collectKeys(result, leftEnv)

	var content string
	var err error

	switch opts.Format {
	case FormatJSON:
		content, err = renderJSON(keys, leftEnv, opts.IncludeVals)
	case FormatMarkdown:
		content = renderMarkdown(keys, leftEnv, opts.IncludeVals)
	default:
		content = renderDotEnv(keys, leftEnv, opts.IncludeVals)
	}

	if err != nil {
		return fmt.Errorf("exporter: render failed: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("exporter: write failed: %w", err)
	}
	return nil
}

func collectKeys(result *differ.Result, base map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range base {
		seen[k] = struct{}{}
	}
	for _, d := range result.Missing {
		seen[d.Key] = struct{}{}
	}
	for _, d := range result.Mismatched {
		seen[d.Key] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func renderDotEnv(keys []string, env map[string]string, includeVals bool) string {
	var sb strings.Builder
	for _, k := range keys {
		if includeVals {
			sb.WriteString(fmt.Sprintf("%s=%s\n", k, env[k]))
		} else {
			sb.WriteString(fmt.Sprintf("%s=\n", k))
		}
	}
	return sb.String()
}

func renderJSON(keys []string, env map[string]string, includeVals bool) (string, error) {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		if includeVals {
			m[k] = env[k]
		} else {
			m[k] = ""
		}
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

func renderMarkdown(keys []string, env map[string]string, includeVals bool) string {
	var sb strings.Builder
	sb.WriteString("| Key | Value |\n")
	sb.WriteString("|-----|-------|\n")
	for _, k := range keys {
		val := ""
		if includeVals {
			val = env[k]
		}
		sb.WriteString(fmt.Sprintf("| %s | %s |\n", k, val))
	}
	return sb.String()
}
