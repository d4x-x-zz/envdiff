// Package freezer provides functionality to lock (freeze) an env map,
// producing a read-only snapshot with change-detection helpers.
package freezer

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// DefaultOptions returns a safe default Options.
func DefaultOptions() Options {
	return Options{
		IncludeValues: true,
	}
}

// Options controls Freeze behaviour.
type Options struct {
	// IncludeValues includes key values in the frozen hash.
	// Set to false to produce a keys-only fingerprint.
	IncludeValues bool

	// KeyPrefix, when non-empty, restricts the freeze to keys with that prefix.
	KeyPrefix string
}

// Frozen is the result of freezing an env map.
type Frozen struct {
	Keys      []string          // sorted list of frozen keys
	Values    map[string]string // snapshot of key→value pairs
	Fingerprint string          // sha256 hex digest of the frozen state
}

// Freeze captures the given env map into a Frozen snapshot.
func Freeze(env map[string]string, opts Options) *Frozen {
	values := make(map[string]string)
	for k, v := range env {
		if opts.KeyPrefix != "" && !strings.HasPrefix(k, opts.KeyPrefix) {
			continue
		}
		values[k] = v
	}

	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		if opts.IncludeValues {
			fmt.Fprintf(h, "%s=%s\n", k, values[k])
		} else {
			fmt.Fprintf(h, "%s\n", k)
		}
	}

	return &Frozen{
		Keys:        keys,
		Values:      values,
		Fingerprint: fmt.Sprintf("%x", h.Sum(nil)),
	}
}

// Changed returns true when the current env no longer matches the frozen snapshot.
func (f *Frozen) Changed(current map[string]string, opts Options) bool {
	new := Freeze(current, opts)
	return new.Fingerprint != f.Fingerprint
}
