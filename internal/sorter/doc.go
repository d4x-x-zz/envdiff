// Package sorter provides key-sorting utilities for env maps.
//
// It supports two modes:
//
//   - Alphabetical: sorts all keys from A to Z.
//   - GroupByPrefix: clusters keys that share a common prefix
//     (the segment before the first underscore), then sorts
//     groups and keys within each group alphabetically.
//
// Example usage:
//
//	env := map[string]string{
//		"DB_HOST": "localhost",
//		"DB_PORT": "5432",
//		"APP_ENV": "production",
//	}
//
//	keys := sorter.Sort(env, sorter.Options{GroupByPrefix: true})
//	// keys => ["APP_ENV", "DB_HOST", "DB_PORT"]
package sorter
