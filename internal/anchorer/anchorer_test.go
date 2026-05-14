package anchorer

import (
	"testing"
)

func TestAnchor_NoOptions_NoChanges(t *testing.T) {
	anchor := map[string]string{"A": "1", "B": "2"}
	target := map[string]string{"A": "1", "B": "2"}
	opts := Options{FillMissing: false, RemoveExtra: false}

	r := Anchor(anchor, target, opts)

	if len(r.Added) != 0 || len(r.Removed) != 0 || len(r.Extra) != 0 {
		t.Fatalf("expected clean result, got added=%v removed=%v extra=%v", r.Added, r.Removed, r.Extra)
	}
	if r.Output["A"] != "1" || r.Output["B"] != "2" {
		t.Fatal("output values changed unexpectedly")
	}
}

func TestAnchor_FillMissing(t *testing.T) {
	anchor := map[string]string{"A": "1", "B": "2", "C": "3"}
	target := map[string]string{"A": "1"}
	opts := Options{FillMissing: true, DefaultValue: "CHANGE_ME"}

	r := Anchor(anchor, target, opts)

	if len(r.Added) != 2 {
		t.Fatalf("expected 2 added keys, got %v", r.Added)
	}
	if r.Output["B"] != "CHANGE_ME" || r.Output["C"] != "CHANGE_ME" {
		t.Fatal("missing keys not filled with default value")
	}
}

func TestAnchor_RemoveExtra(t *testing.T) {
	anchor := map[string]string{"A": "1"}
	target := map[string]string{"A": "1", "X": "extra", "Y": "also extra"}
	opts := Options{RemoveExtra: true}

	r := Anchor(anchor, target, opts)

	if len(r.Removed) != 2 {
		t.Fatalf("expected 2 removed keys, got %v", r.Removed)
	}
	if _, ok := r.Output["X"]; ok {
		t.Fatal("extra key X should have been removed")
	}
}

func TestAnchor_ExtraReportedButNotRemoved(t *testing.T) {
	anchor := map[string]string{"A": "1"}
	target := map[string]string{"A": "1", "Z": "extra"}
	opts := Options{RemoveExtra: false}

	r := Anchor(anchor, target, opts)

	if len(r.Extra) != 1 || r.Extra[0] != "Z" {
		t.Fatalf("expected extra=[Z], got %v", r.Extra)
	}
	if len(r.Removed) != 0 {
		t.Fatalf("expected nothing removed, got %v", r.Removed)
	}
	if _, ok := r.Output["Z"]; !ok {
		t.Fatal("extra key Z should still be present in output")
	}
}

func TestAnchor_DefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if !opts.FillMissing {
		t.Fatal("FillMissing should default to true")
	}
	if opts.RemoveExtra {
		t.Fatal("RemoveExtra should default to false")
	}
	if opts.DefaultValue != "" {
		t.Fatalf("DefaultValue should be empty string, got %q", opts.DefaultValue)
	}
}

func TestAnchor_EmptyTarget(t *testing.T) {
	anchor := map[string]string{"A": "1", "B": "2"}
	target := map[string]string{}
	opts := DefaultOptions()

	r := Anchor(anchor, target, opts)

	if len(r.Added) != 2 {
		t.Fatalf("expected 2 added keys, got %v", r.Added)
	}
	if len(r.Output) != 2 {
		t.Fatalf("expected output to have 2 keys, got %d", len(r.Output))
	}
}
