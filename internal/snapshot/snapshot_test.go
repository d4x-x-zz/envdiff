package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/snapshot"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := snapshot.Save(path, "test", env); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if s.Label != "test" {
		t.Errorf("label: got %q want %q", s.Label, "test")
	}
	if s.Env["FOO"] != "bar" {
		t.Errorf("FOO: got %q want %q", s.Env["FOO"], "bar")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCompare_Clean(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	old := &snapshot.Snapshot{Env: env}
	new := &snapshot.Snapshot{Env: map[string]string{"A": "1", "B": "2"}}
	res := snapshot.Compare(old, new)
	if !res.Clean() {
		t.Error("expected clean result")
	}
}

func TestCompare_Added(t *testing.T) {
	old := &snapshot.Snapshot{Env: map[string]string{"A": "1"}}
	new := &snapshot.Snapshot{Env: map[string]string{"A": "1", "B": "2"}}
	res := snapshot.Compare(old, new)
	if len(res.Added) != 1 || res.Added[0] != "B" {
		t.Errorf("Added: got %v", res.Added)
	}
}

func TestCompare_Removed(t *testing.T) {
	old := &snapshot.Snapshot{Env: map[string]string{"A": "1", "B": "2"}}
	new := &snapshot.Snapshot{Env: map[string]string{"A": "1"}}
	res := snapshot.Compare(old, new)
	if len(res.Removed) != 1 || res.Removed[0] != "B" {
		t.Errorf("Removed: got %v", res.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	old := &snapshot.Snapshot{Env: map[string]string{"A": "old"}}
	new := &snapshot.Snapshot{Env: map[string]string{"A": "new"}}
	res := snapshot.Compare(old, new)
	if len(res.Changed) != 1 || res.Changed[0] != "A" {
		t.Errorf("Changed: got %v", res.Changed)
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := snapshot.Save("/no/such/dir/snap.json", "x", map[string]string{})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func init() {
	_ = os.Getenv // ensure os is used
}
