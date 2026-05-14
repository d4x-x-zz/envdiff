// Package digester computes a deterministic hash digest of an env map,
// useful for detecting changes between two snapshots or pipeline stages.
package digester

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// Options controls how the digest is computed.
type Options struct {
	// IncludeValues includes key values in the hash. When false, only key
	// names are hashed (useful for structural comparison).
	IncludeValues bool
	// KeyPrefix restricts hashing to keys that start with the given prefix.
	KeyPrefix string
}

// DefaultOptions returns sensible defaults: values are included, no prefix filter.
func DefaultOptions() Options {
	return Options{
		IncludeValues: true,
	}
}

// Result holds the output of a Digest call.
type Result struct {
	// Hex is the hex-encoded SHA-256 digest.
	Hex string
	// KeyCount is the number of keys that contributed to the digest.
	KeyCount int
}

// Digest computes a deterministic SHA-256 hash of the provided env map.
// Keys are sorted before hashing to ensure consistency regardless of map
// iteration order.
func Digest(env map[string]string, opts Options) Result {
	keys := make([]string, 0, len(env))
	for k := range env {
		if opts.KeyPrefix != "" && !strings.HasPrefix(k, opts.KeyPrefix) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		if opts.IncludeValues {
			fmt.Fprintf(h, "%s=%s\n", k, env[k])
		} else {
			fmt.Fprintf(h, "%s\n", k)
		}
	}

	return Result{
		Hex:      fmt.Sprintf("%x", h.Sum(nil)),
		KeyCount: len(keys),
	}
}
