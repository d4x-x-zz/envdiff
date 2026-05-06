package cli

import (
	"flag"
	"fmt"
	"os"

	"envdiff/internal/parser"
	"envdiff/internal/templater"
)

type templateArgs struct {
	input        string
	output       string
	typed        bool
	commentOrig  bool
}

func parseTemplateArgs(args []string) (templateArgs, error) {
	fs := flag.NewFlagSet("template", flag.ContinueOnError)
	output := fs.String("output", "", "write template to file instead of stdout")
	typed := fs.Bool("typed", true, "use typed placeholders (<string>, <number>, <bool>)")
	comment := fs.Bool("comment-original", false, "include original value as comment")

	if err := fs.Parse(args); err != nil {
		return templateArgs{}, err
	}
	if fs.NArg() < 1 {
		return templateArgs{}, fmt.Errorf("usage: envdiff template [flags] <file>")
	}
	return templateArgs{
		input:       fs.Arg(0),
		output:      *output,
		typed:       *typed,
		commentOrig: *comment,
	}, nil
}

// RunTemplate generates a .env.template from an existing .env file.
func RunTemplate(args []string, stdout *os.File) error {
	tArgs, err := parseTemplateArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(tArgs.input)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", tArgs.input, err)
	}

	opts := templater.Options{
		UseTypedPlaceholders: tArgs.typed,
		CommentOriginal:      tArgs.commentOrig,
	}
	result := templater.Generate(env, opts)

	if tArgs.output != "" {
		if err := os.WriteFile(tArgs.output, []byte(result), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(stdout, "template written to %s\n", tArgs.output)
		return nil
	}

	fmt.Fprint(stdout, result)
	return nil
}
