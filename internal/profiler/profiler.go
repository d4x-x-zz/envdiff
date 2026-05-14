// Package profiler analyses an env map and produces a usage profile:
// key count, value length stats, type distribution, and density metrics.
package profiler

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

// Options controls which analyses are performed.
type Options struct {
	IncludeTypeBreakdown bool
	IncludeLengthStats   bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		IncludeTypeBreakdown: true,
		IncludeLengthStats:   true,
	}
}

// Profile holds the analysis results for an env map.
type Profile struct {
	TotalKeys     int
	EmptyValues   int
	TypeBreakdown map[string]int // "bool", "int", "float", "url", "string"
	AvgValueLen   float64
	MaxValueLen   int
	MinValueLen   int
	Density       float64 // fraction of non-empty values
}

// Analyze profiles the given env map according to opts.
func Analyze(env map[string]string, opts Options) Profile {
	p := Profile{
		TotalKeys:     len(env),
		TypeBreakdown: map[string]int{},
		MinValueLen:   math.MaxInt32,
	}

	if len(env) == 0 {
		p.MinValueLen = 0
		return p
	}

	totalLen := 0
	for _, v := range env {
		if v == "" {
			p.EmptyValues++
		}
		l := len(v)
		totalLen += l
		if opts.IncludeLengthStats {
			if l > p.MaxValueLen {
				p.MaxValueLen = l
			}
			if l < p.MinValueLen {
				p.MinValueLen = l
			}
		}
		if opts.IncludeTypeBreakdown {
			p.TypeBreakdown[inferType(v)]++
		}
	}

	if opts.IncludeLengthStats {
		p.AvgValueLen = float64(totalLen) / float64(len(env))
	}
	nonEmpty := len(env) - p.EmptyValues
	p.Density = float64(nonEmpty) / float64(len(env))
	return p
}

// SortedTypes returns type names sorted alphabetically.
func SortedTypes(p Profile) []string {
	keys := make([]string, 0, len(p.TypeBreakdown))
	for k := range p.TypeBreakdown {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func inferType(v string) string {
	if v == "" {
		return "string"
	}
	lv := strings.ToLower(v)
	if lv == "true" || lv == "false" {
		return "bool"
	}
	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return "int"
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return "float"
	}
	if strings.HasPrefix(lv, "http://") || strings.HasPrefix(lv, "https://") {
		return "url"
	}
	return "string"
}
