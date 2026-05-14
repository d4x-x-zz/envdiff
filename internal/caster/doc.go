// Package caster infers the likely data type of each value in a .env map.
//
// Supported types:
//
//   - bool   — "true", "false", "1", "0"
//   - int    — any integer literal
//   - float  — any floating-point literal
//   - url    — values starting with http:// or https://
//   - dsn    — connection-string style values (scheme://user@host/db)
//   - string — everything else
//
// Example:
//
//	results := caster.Cast(env, caster.DefaultOptions())
//	for _, r := range results {
//		fmt.Printf("%s = %s (%s)\n", r.Key, r.Value, r.Type)
//	}
package caster
