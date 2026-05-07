// Package snapshot provides functionality to save and compare .env snapshots over time.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of an env map at a point in time.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label"`
	Env       map[string]string `json:"env"`
}

// Save writes a snapshot of the given env map to the specified file path.
func Save(path, label string, env map[string]string) error {
	s := Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Env:       env,
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// Compare returns keys added, removed, or changed between two snapshots.
func Compare(old, new *Snapshot) CompareResult {
	res := CompareResult{}
	for k, v := range new.Env {
		oldVal, exists := old.Env[k]
		if !exists {
			res.Added = append(res.Added, k)
		} else if oldVal != v {
			res.Changed = append(res.Changed, k)
		}
	}
	for k := range old.Env {
		if _, exists := new.Env[k]; !exists {
			res.Removed = append(res.Removed, k)
		}
	}
	return res
}

// CompareResult holds the diff between two snapshots.
type CompareResult struct {
	Added   []string
	Removed []string
	Changed []string
}

// Clean returns true when there are no differences.
func (r CompareResult) Clean() bool {
	return len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Changed) == 0
}
