package validator_test

import (
	"testing"

	"github.com/user/envdiff/internal/validator"
)

func TestValidate_EmptyValue(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "",
		"DB_PORT": "5432",
	}
	opts := validator.Options{WarnEmpty: true}
	issues := validator.Validate(env, opts)

	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "DB_HOST" {
		t.Errorf("expected issue for DB_HOST, got %s", issues[0].Key)
	}
}

func TestValidate_PlaceholderValue(t *testing.T) {
	env := map[string]string{
		"API_KEY":  "changeme",
		"APP_NAME": "myapp",
	}
	opts := validator.Options{WarnPlaceholder: true}
	issues := validator.Validate(env, opts)

	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "API_KEY" {
		t.Errorf("expected issue for API_KEY, got %s", issues[0].Key)
	}
}

func TestValidate_NamingConvention(t *testing.T) {
	env := map[string]string{
		"db_host": "localhost",
		"DB_PORT": "5432",
	}
	opts := validator.Options{WarnNaming: true}
	issues := validator.Validate(env, opts)

	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "db_host" {
		t.Errorf("expected issue for db_host, got %s", issues[0].Key)
	}
}

func TestValidate_NoIssues(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	issues := validator.Validate(env, validator.DefaultOptions())
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestValidate_MultipleChecks(t *testing.T) {
	env := map[string]string{
		"SECRET": "",
		"TOKEN":  "your_token_here",
	}
	opts := validator.Options{WarnEmpty: true, WarnPlaceholder: true}
	issues := validator.Validate(env, opts)

	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
}

func TestIssue_String(t *testing.T) {
	i := validator.Issue{Key: "FOO", Message: "value is empty"}
	if i.String() != "FOO: value is empty" {
		t.Errorf("unexpected string: %s", i.String())
	}
}
