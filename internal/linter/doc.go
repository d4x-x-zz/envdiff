// Package linter provides style and convention checks for .env key names.
//
// It is intentionally separate from the validator package: the validator
// focuses on value quality (empty values, placeholders, etc.) while the
// linter focuses on key naming conventions and file-level patterns.
//
// Basic usage:
//
//	env, _ := parser.ParseFile("production.env")
//	issues := linter.Lint(env, linter.DefaultOptions())
//	for _, iss := range issues {
//		fmt.Printf("[%s] %s\n", iss.Key, iss.Message)
//	}
//
// Available checks (all toggleable via Options):
//   - CheckUpperCase    – keys must be UPPER_CASE
//   - CheckNoSpaces     – keys must not contain whitespace
//   - CheckNoLeadDigit  – keys must not start with a digit
//   - CheckNoDupPrefix  – warn when one prefix dominates the file
package linter
