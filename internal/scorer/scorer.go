// Package scorer computes a quality score for a .env map based on
// configurable checks such as key naming, value completeness, and
// absence of placeholder values.
package scorer

import (
	"strings"
)

// Options controls which checks contribute to the score.
type Options struct {
	// MaxScore is the ceiling score (default 100).
	MaxScore int
	// PenalizeEmpty deducts points for keys with empty values.
	PenalizeEmpty bool
	// PenalizePlaceholders deducts points for placeholder values like CHANGE_ME.
	PenalizePlaceholders bool
	// PenalizeLowercase deducts points for keys that are not fully uppercase.
	PenalizeLowercase bool
	// PenalizeNoValue deducts points for keys whose value equals the key name.
	PenalizeNoValue bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxScore:             100,
		PenalizeEmpty:        true,
		PenalizePlaceholders: true,
		PenalizeLowercase:    true,
		PenalizeNoValue:      true,
	}
}

// Result holds the computed score and a breakdown of deductions.
type Result struct {
	Score      int
	MaxScore   int
	Deductions []string
}

// Grade returns a letter grade for the score.
func (r Result) Grade() string {
	pct := 0
	if r.MaxScore > 0 {
		pct = r.Score * 100 / r.MaxScore
	}
	switch {
	case pct >= 90:
		return "A"
	case pct >= 75:
		return "B"
	case pct >= 60:
		return "C"
	case pct >= 40:
		return "D"
	default:
		return "F"
	}
}

var placeholders = []string{"change_me", "changeme", "todo", "fixme", "placeholder", "your_", "<", "example"}

// Score evaluates env and returns a Result.
func Score(env map[string]string, opts Options) Result {
	if opts.MaxScore <= 0 {
		opts.MaxScore = 100
	}
	if len(env) == 0 {
		return Result{Score: opts.MaxScore, MaxScore: opts.MaxScore}
	}

	total := len(env)
	deductions := []string{}
	penalty := 0

	for k, v := range env {
		low := strings.ToLower(v)
		if opts.PenalizeEmpty && v == "" {
			penalty++
			deductions = append(deductions, k+": empty value")
			continue
		}
		if opts.PenalizePlaceholders {
			for _, ph := range placeholders {
				if strings.Contains(low, ph) {
					penalty++
					deductions = append(deductions, k+": placeholder value")
					break
				}
			}
		}
		if opts.PenalizeLowercase && k != strings.ToUpper(k) {
			penalty++
			deductions = append(deductions, k+": key not uppercase")
		}
		if opts.PenalizeNoValue && strings.EqualFold(k, v) {
			penalty++
			deductions = append(deductions, k+": value mirrors key")
		}
	}

	score := opts.MaxScore - (penalty*opts.MaxScore)/total
	if score < 0 {
		score = 0
	}
	return Result{Score: score, MaxScore: opts.MaxScore, Deductions: deductions}
}
