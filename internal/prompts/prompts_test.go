package prompts_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/prompts"
)

func opts(in string) prompts.Options {
	o := prompts.DefaultOptions()
	o.In = strings.NewReader(in)
	o.Out = &bytes.Buffer{}
	return o
}

func TestFill_SkipsFilledValues(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	o := opts("")
	o.SkipFilled = true
	res, err := prompts.Fill(env, o)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["HOST"] != "localhost" || res["PORT"] != "8080" {
		t.Errorf("expected original values preserved, got %v", res)
	}
}

func TestFill_PromptsForEmptyValue(t *testing.T) {
	env := map[string]string{"SECRET": ""}
	o := opts("mysecret\n")
	res, err := prompts.Fill(env, o)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["SECRET"] != "mysecret" {
		t.Errorf("expected 'mysecret', got %q", res["SECRET"])
	}
}

func TestFill_PromptsForPlaceholder(t *testing.T) {
	env := map[string]string{"API_KEY": "<YOUR_API_KEY>"}
	o := opts("abc123\n")
	res, err := prompts.Fill(env, o)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["API_KEY"] != "abc123" {
		t.Errorf("expected 'abc123', got %q", res["API_KEY"])
	}
}

func TestFill_KeepsOriginalOnBlankInput(t *testing.T) {
	env := map[string]string{"TOKEN": "<CHANGE_ME>"}
	o := opts("\n") // user presses enter without typing
	res, err := prompts.Fill(env, o)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["TOKEN"] != "<CHANGE_ME>" {
		t.Errorf("expected placeholder preserved, got %q", res["TOKEN"])
	}
}

func TestFill_MultipleKeys_OrderedPrompts(t *testing.T) {
	env := map[string]string{"B_KEY": "", "A_KEY": ""}
	out := &bytes.Buffer{}
	o := prompts.DefaultOptions()
	o.In = strings.NewReader("val_a\nval_b\n")
	o.Out = out
	res, err := prompts.Fill(env, o)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["A_KEY"] != "val_a" {
		t.Errorf("A_KEY: expected val_a, got %q", res["A_KEY"])
	}
	if res["B_KEY"] != "val_b" {
		t.Errorf("B_KEY: expected val_b, got %q", res["B_KEY"])
	}
}

func TestDefaultOptions(t *testing.T) {
	o := prompts.DefaultOptions()
	if !o.SkipFilled {
		t.Error("expected SkipFilled to be true by default")
	}
	if len(o.PlaceholderMarkers) == 0 {
		t.Error("expected at least one placeholder marker")
	}
}
