package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/freezer"
	"github.com/user/envdiff/internal/parser"
)

type freezeArgs struct {
	file      string
	keyPrefix string
	keysOnly  bool
	format    string
	output    string
}

func parseFreezeArgs(args []string) (freezeArgs, error) {
	fs := flag.NewFlagSet("freeze", flag.ContinueOnError)
	prefix := fs.String("prefix", "", "only freeze keys with this prefix")
	keysOnly := fs.Bool("keys-only", false, "fingerprint keys only, ignore values")
	format := fs.String("format", "text", "output format: text|json")
	out := fs.String("output", "", "write result to file instead of stdout")
	if err := fs.Parse(args); err != nil {
		return freezeArgs{}, err
	}
	if fs.NArg() < 1 {
		return freezeArgs{}, fmt.Errorf("freeze: requires a .env file argument")
	}
	return freezeArgs{
		file:      fs.Arg(0),
		keyPrefix: *prefix,
		keysOnly:  *keysOnly,
		format:    *format,
		output:    *out,
	}, nil
}

// RunFreeze is the entry point for the freeze sub-command.
func RunFreeze(args []string, stdout io.Writer) error {
	a, err := parseFreezeArgs(args)
	if err != nil {
		return err
	}
	env, err := parser.ParseFile(a.file)
	if err != nil {
		return fmt.Errorf("freeze: %w", err)
	}
	opts := freezer.Options{
		IncludeValues: !a.keysOnly,
		KeyPrefix:     a.keyPrefix,
	}
	f := freezer.Freeze(env, opts)

	w := stdout
	if a.output != "" {
		file, err := os.Create(a.output)
		if err != nil {
			return fmt.Errorf("freeze: %w", err)
		}
		defer file.Close()
		w = file
	}

	switch a.format {
	case "json":
		return json.NewEncoder(w).Encode(map[string]interface{}{
			"fingerprint": f.Fingerprint,
			"key_count":   len(f.Keys),
			"keys":        f.Keys,
		})
	default:
		fmt.Fprintf(w, "fingerprint: %s\n", f.Fingerprint)
		fmt.Fprintf(w, "keys:        %d\n", len(f.Keys))
		for _, k := range f.Keys {
			fmt.Fprintf(w, "  %s\n", k)
		}
	}
	return nil
}
