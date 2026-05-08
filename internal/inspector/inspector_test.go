package inspector

import (
	"testing"
)

func TestInspect_TotalKeys(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	s := Inspect(env, DefaultOptions())
	if s.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", s.TotalKeys)
	}
}

func TestInspect_EmptyValues(t *testing.T) {
	env := map[string]string{"A": "", "B": "hello", "C": ""}
	s := Inspect(env, DefaultOptions())
	if s.EmptyValues != 2 {
		t.Errorf("expected 2 empty values, got %d", s.EmptyValues)
	}
}

func TestInspect_SensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret",
		"API_TOKEN":   "tok",
		"APP_NAME":    "myapp",
	}
	s := Inspect(env, DefaultOptions())
	if len(s.SensitiveKeys) != 2 {
		t.Errorf("expected 2 sensitive keys, got %d: %v", len(s.SensitiveKeys), s.SensitiveKeys)
	}
}

func TestInspect_Placeholders(t *testing.T) {
	env := map[string]string{
		"A": "<your-value>",
		"B": "CHANGE_ME",
		"C": "real-value",
	}
	s := Inspect(env, DefaultOptions())
	if s.Placeholders != 2 {
		t.Errorf("expected 2 placeholders, got %d", s.Placeholders)
	}
}

func TestInspect_UniqueValues(t *testing.T) {
	env := map[string]string{"A": "foo", "B": "foo", "C": "bar", "D": ""}
	s := Inspect(env, DefaultOptions())
	// empty values are excluded from unique count
	if s.UniqueValues != 2 {
		t.Errorf("expected 2 unique values, got %d", s.UniqueValues)
	}
}

func TestInspect_EmptyMap(t *testing.T) {
	s := Inspect(map[string]string{}, DefaultOptions())
	if s.TotalKeys != 0 || s.EmptyValues != 0 || s.Placeholders != 0 {
		t.Errorf("expected all zero for empty map, got %+v", s)
	}
}

func TestInspect_CustomOptions(t *testing.T) {
	opts := Options{
		SensitivePatterns:   []string{"PRIVATE"},
		PlaceholderPrefixes: []string{"TODO_"},
	}
	env := map[string]string{
		"PRIVATE_KEY": "abc",
		"API_TOKEN":   "tok",
		"VAL":         "TODO_FILL",
	}
	s := Inspect(env, opts)
	if len(s.SensitiveKeys) != 1 || s.SensitiveKeys[0] != "PRIVATE_KEY" {
		t.Errorf("unexpected sensitive keys: %v", s.SensitiveKeys)
	}
	if s.Placeholders != 1 {
		t.Errorf("expected 1 placeholder, got %d", s.Placeholders)
	}
}
