package linter

import (
	"testing"
)

func TestLint_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	issues := Lint(env, DefaultOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_LowerCaseKey(t *testing.T) {
	env := map[string]string{"app_host": "localhost"}
	issues := Lint(env, DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "app_host" {
		t.Errorf("unexpected key %q", issues[0].Key)
	}
}

func TestLint_KeyWithSpace(t *testing.T) {
	env := map[string]string{"APP HOST": "val"}
	opts := DefaultOptions()
	opts.CheckUpperCase = false // isolate the space check
	issues := Lint(env, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Message != "key contains whitespace" {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestLint_KeyStartsWithDigit(t *testing.T) {
	env := map[string]string{"1APP": "val"}
	opts := DefaultOptions()
	opts.CheckUpperCase = false
	issues := Lint(env, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Message != "key starts with a digit" {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestLint_DupPrefix(t *testing.T) {
	env := map[string]string{
		"APP_A": "1",
		"APP_B": "2",
		"APP_C": "3",
		"OTHER": "x",
	}
	opts := Options{CheckNoDupPrefix: true}
	issues := Lint(env, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 dup-prefix issue, got %d: %v", len(issues), issues)
	}
	if issues[0].Key != "APP_*" {
		t.Errorf("expected key APP_*, got %q", issues[0].Key)
	}
}

func TestLint_EmptyMap(t *testing.T) {
	issues := Lint(map[string]string{}, DefaultOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues on empty map, got %d", len(issues))
	}
}

func TestLint_DisabledChecks(t *testing.T) {
	env := map[string]string{"lower key": "val"}
	opts := Options{} // all checks disabled
	issues := Lint(env, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues with all checks off, got %d", len(issues))
	}
}
