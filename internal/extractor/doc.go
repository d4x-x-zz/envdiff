// Package extractor provides utilities for pulling a selective subset of
// key-value pairs from an env map.
//
// Keys can be selected either by an explicit allow-list or by one or more
// glob patterns (e.g. "DB_*", "APP_*"). When both are supplied, the explicit
// key list takes precedence.
//
// Example usage:
//
//	opts := extractor.DefaultOptions()
//	opts.Patterns = []string{"DB_*"}
//	result := extractor.Extract(envMap, opts)
//	fmt.Println(result.Env)    // extracted keys
//	fmt.Println(result.Missed) // keys requested but absent
package extractor
