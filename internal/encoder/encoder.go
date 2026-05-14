// Package encoder converts an env map into various serialization formats
// such as shell export statements, Docker --env-file format, and YAML.
package encoder

import (
	"fmt"
	"sort"
	"strings"
)

// Format represents the target encoding format.
type Format string

const (
	FormatShell  Format = "shell"
	FormatDocker Format = "docker"
	FormatYAML   Format = "yaml"
)

// Options controls encoder behaviour.
type Options struct {
	Format    Format
	SortKeys  bool
	OmitEmpty bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Format:    FormatShell,
		SortKeys:  true,
		OmitEmpty: false,
	}
}

// Encode serializes env into the requested format.
func Encode(env map[string]string, opts Options) (string, error) {
	keys := collectKeys(env, opts)

	switch opts.Format {
	case FormatShell:
		return renderShell(env, keys), nil
	case FormatDocker:
		return renderDocker(env, keys), nil
	case FormatYAML:
		return renderYAML(env, keys), nil
	default:
		return "", fmt.Errorf("unknown format %q", opts.Format)
	}
}

func collectKeys(env map[string]string, opts Options) []string {
	keys := make([]string, 0, len(env))
	for k, v := range env {
		if opts.OmitEmpty && v == "" {
			continue
		}
		keys = append(keys, k)
	}
	if opts.SortKeys {
		sort.Strings(keys)
	}
	return keys
}

func renderShell(env map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String()
}

func renderDocker(env map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}
	return sb.String()
}

func renderYAML(env map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		if strings.ContainsAny(v, ":\"'#") || v == "" {
			fmt.Fprintf(&sb, "%s: %q\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s: %s\n", k, v)
		}
	}
	return sb.String()
}
