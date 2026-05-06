package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envdiff/internal/differ"
)

func TestReporter_NoDifferences_Text(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)
	if err := r.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-differences message, got: %q", buf.String())
	}
}

func TestReporter_MissingRight_Text(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)
	results := []differ.Result{
		{Key: "DB_HOST", Kind: differ.MissingInRight, LeftVal: "localhost"},
	}
	r.Write(results)
	out := buf.String()
	if !strings.Contains(out, "MISSING_RIGHT") || !strings.Contains(out, "DB_HOST") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestReporter_Mismatch_Text(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)
	results := []differ.Result{
		{Key: "PORT", Kind: differ.ValueMismatch, LeftVal: "3000", RightVal: "4000"},
	}
	r.Write(results)
	out := buf.String()
	if !strings.Contains(out, "MISMATCH") {
		t.Errorf("expected MISMATCH in output, got: %q", out)
	}
	if !strings.Contains(out, "3000") || !strings.Contains(out, "4000") {
		t.Errorf("expected both values in output, got: %q", out)
	}
}

func TestReporter_NoDifferences_JSON(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatJSON)
	r.Write([]differ.Result{})
	out := buf.String()
	if !strings.Contains(out, `"differences":[]`) {
		t.Errorf("expected empty JSON array, got: %q", out)
	}
}

func TestReporter_MissingLeft_JSON(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatJSON)
	results := []differ.Result{
		{Key: "SECRET", Kind: differ.MissingInLeft, RightVal: "abc123"},
	}
	r.Write(results)
	out := buf.String()
	if !strings.Contains(out, `"SECRET"`) {
		t.Errorf("expected key in JSON output, got: %q", out)
	}
	if !strings.Contains(out, `missing_left`) {
		t.Errorf("expected kind in JSON output, got: %q", out)
	}
}

func TestReporter_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)
	results := []differ.Result{
		{Key: "Z_KEY", Kind: differ.MissingInRight},
		{Key: "A_KEY", Kind: differ.MissingInRight},
		{Key: "M_KEY", Kind: differ.MissingInRight},
	}
	r.Write(results)
	out := buf.String()
	aIdx := strings.Index(out, "A_KEY")
	mIdx := strings.Index(out, "M_KEY")
	zIdx := strings.Index(out, "Z_KEY")
	if !(aIdx < mIdx && mIdx < zIdx) {
		t.Errorf("output not sorted alphabetically: %q", out)
	}
}
