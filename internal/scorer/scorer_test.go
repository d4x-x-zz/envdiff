package scorer

import (
	"testing"
)

func TestScore_PerfectEnv(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"SECRET_KEY":   "s3cr3t!",
		"PORT":         "8080",
	}
	r := Score(env, DefaultOptions())
	if r.Score != 100 {
		t.Errorf("expected 100, got %d", r.Score)
	}
	if r.Grade() != "A" {
		t.Errorf("expected grade A, got %s", r.Grade())
	}
}

func TestScore_EmptyValues(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "",
		"SECRET_KEY":   "real-value",
	}
	opts := DefaultOptions()
	r := Score(env, opts)
	if r.Score >= 100 {
		t.Errorf("expected penalty for empty value, got score %d", r.Score)
	}
	if len(r.Deductions) == 0 {
		t.Error("expected at least one deduction")
	}
}

func TestScore_PlaceholderValue(t *testing.T) {
	env := map[string]string{
		"API_KEY": "CHANGE_ME",
		"HOST":    "localhost",
	}
	r := Score(env, DefaultOptions())
	if r.Score >= 100 {
		t.Errorf("expected penalty for placeholder, got %d", r.Score)
	}
	found := false
	for _, d := range r.Deductions {
		if d == "API_KEY: placeholder value" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected placeholder deduction, got %v", r.Deductions)
	}
}

func TestScore_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"db_host": "localhost",
		"PORT":    "5432",
	}
	r := Score(env, DefaultOptions())
	if r.Score >= 100 {
		t.Errorf("expected penalty for lowercase key, got %d", r.Score)
	}
}

func TestScore_ValueMirrorsKey(t *testing.T) {
	env := map[string]string{
		"TOKEN": "token",
	}
	r := Score(env, DefaultOptions())
	if r.Score >= 100 {
		t.Errorf("expected penalty for value mirroring key, got %d", r.Score)
	}
}

func TestScore_EmptyMap(t *testing.T) {
	r := Score(map[string]string{}, DefaultOptions())
	if r.Score != 100 {
		t.Errorf("empty map should yield max score, got %d", r.Score)
	}
}

func TestScore_GradeThresholds(t *testing.T) {
	cases := []struct {
		score    int
		max      int
		expected string
	}{
		{95, 100, "A"},
		{80, 100, "B"},
		{65, 100, "C"},
		{45, 100, "D"},
		{20, 100, "F"},
	}
	for _, tc := range cases {
		r := Result{Score: tc.score, MaxScore: tc.max}
		if g := r.Grade(); g != tc.expected {
			t.Errorf("score %d/%d: expected grade %s, got %s", tc.score, tc.max, tc.expected, g)
		}
	}
}

func TestScore_DisabledChecks(t *testing.T) {
	env := map[string]string{
		"db_host": "",
		"api_key": "CHANGE_ME",
	}
	opts := Options{
		MaxScore:             100,
		PenalizeEmpty:        false,
		PenalizePlaceholders: false,
		PenalizeLowercase:    false,
		PenalizeNoValue:      false,
	}
	r := Score(env, opts)
	if r.Score != 100 {
		t.Errorf("all checks disabled, expected 100, got %d", r.Score)
	}
}
