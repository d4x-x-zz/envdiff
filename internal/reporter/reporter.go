package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/envdiff/internal/differ"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes diff results to an output stream.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a new Reporter writing to the given writer.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out, format: format}
}

// Write outputs the diff results according to the configured format.
func (r *Reporter) Write(results []differ.Result) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(results)
	default:
		return r.writeText(results)
	}
}

func (r *Reporter) writeText(results []differ.Result) error {
	if len(results) == 0 {
		fmt.Fprintln(r.out, "✓ No differences found.")
		return nil
	}

	sorted := make([]differ.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, res := range sorted {
		switch res.Kind {
		case differ.MissingInRight:
			fmt.Fprintf(r.out, "  MISSING_RIGHT  %s\n", res.Key)
		case differ.MissingInLeft:
			fmt.Fprintf(r.out, "  MISSING_LEFT   %s\n", res.Key)
		case differ.ValueMismatch:
			fmt.Fprintf(r.out, "  MISMATCH       %s  (%q != %q)\n", res.Key, res.LeftVal, res.RightVal)
		}
	}
	return nil
}

func (r *Reporter) writeJSON(results []differ.Result) error {
	if len(results) == 0 {
		fmt.Fprintln(r.out, `{"differences":[]}`)
		return nil
	}

	fmt.Fprintln(r.out, `{"differences":[`)
	for i, res := range results {
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Fprintf(r.out, `  {"key":%q,"kind":%q,"left":%q,"right":%q}%s\n`,
			res.Key, res.Kind, res.LeftVal, res.RightVal, comma)
	}
	fmt.Fprintln(r.out, `]}`)
	return nil
}
