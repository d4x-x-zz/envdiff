// Package comparator provides multi-file environment comparison,
// summarising which keys are present, missing, or divergent across
// an arbitrary number of .env files.
package comparator

import "sort"

// FileEnv maps a file label to its parsed key/value pairs.
type FileEnv map[string]map[string]string

// KeyStatus describes how a single key appears across all files.
type KeyStatus struct {
	Key      string
	Values   map[string]string // file label → value
	Missing  []string          // file labels where the key is absent
	Uniform  bool              // true when all present values are identical
}

// Report is the result of comparing multiple env files.
type Report struct {
	Files    []string    // ordered list of file labels
	Statuses []KeyStatus // one entry per unique key
}

// TotalMissing returns the number of keys that are absent in at least one file.
func (r *Report) TotalMissing() int {
	count := 0
	for _, s := range r.Statuses {
		if len(s.Missing) > 0 {
			count++
		}
	}
	return count
}

// TotalMismatched returns the number of keys whose values differ across files.
func (r *Report) TotalMismatched() int {
	count := 0
	for _, s := range r.Statuses {
		if !s.Uniform && len(s.Missing) == 0 {
			count++
		}
	}
	return count
}

// Compare builds a Report from the supplied FileEnv.
func Compare(files FileEnv) *Report {
	// collect ordered file labels
	labels := make([]string, 0, len(files))
	for l := range files {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	// collect all unique keys
	keySet := map[string]struct{}{}
	for _, env := range files {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	statuses := make([]KeyStatus, 0, len(keys))
	for _, k := range keys {
		values := map[string]string{}
		var missing []string
		for _, l := range labels {
			v, ok := files[l][k]
			if !ok {
				missing = append(missing, l)
			} else {
				values[l] = v
			}
		}
		uniform := isUniform(values)
		statuses = append(statuses, KeyStatus{
			Key:     k,
			Values:  values,
			Missing: missing,
			Uniform: uniform,
		})
	}

	return &Report{Files: labels, Statuses: statuses}
}

func isUniform(values map[string]string) bool {
	var ref string
	first := true
	for _, v := range values {
		if first {
			ref = v
			first = false
			continue
		}
		if v != ref {
			return false
		}
	}
	return true
}
