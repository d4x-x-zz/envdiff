package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/snapshot"
)

type snapshotArgs struct {
	envFile  string
	snapFile string
	label    string
	compare  bool
}

func parseSnapshotArgs(args []string) (snapshotArgs, error) {
	fs := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	env := fs.String("env", ".env", "path to .env file")
	snap := fs.String("out", ".env.snap", "path to snapshot file")
	label := fs.String("label", "", "label for this snapshot")
	cmp := fs.Bool("compare", false, "compare env file against existing snapshot")
	if err := fs.Parse(args); err != nil {
		return snapshotArgs{}, err
	}
	return snapshotArgs{envFile: *env, snapFile: *snap, label: *label, compare: *cmp}, nil
}

// RunSnapshot saves or compares a .env snapshot.
func RunSnapshot(args []string, out *os.File) error {
	a, err := parseSnapshotArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(a.envFile)
	if err != nil {
		return fmt.Errorf("snapshot: parse env: %w", err)
	}

	if a.compare {
		oldSnap, err := snapshot.Load(a.snapFile)
		if err != nil {
			return fmt.Errorf("snapshot: load: %w", err)
		}
		newSnap := &snapshot.Snapshot{Env: env}
		res := snapshot.Compare(oldSnap, newSnap)
		if res.Clean() {
			fmt.Fprintln(out, "no changes since snapshot")
			return nil
		}
		if len(res.Added) > 0 {
			fmt.Fprintf(out, "added:   %s\n", strings.Join(res.Added, ", "))
		}
		if len(res.Removed) > 0 {
			fmt.Fprintf(out, "removed: %s\n", strings.Join(res.Removed, ", "))
		}
		if len(res.Changed) > 0 {
			fmt.Fprintf(out, "changed: %s\n", strings.Join(res.Changed, ", "))
		}
		return nil
	}

	if err := snapshot.Save(a.snapFile, a.label, env); err != nil {
		return err
	}
	fmt.Fprintf(out, "snapshot saved to %s\n", a.snapFile)
	return nil
}
