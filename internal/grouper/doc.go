// Package grouper partitions a flat env map into named groups based on key
// prefix conventions.
//
// Keys are split on a configurable separator (default "_"). The portion before
// the first separator becomes the group name. Keys with no separator — or whose
// prefix is not in the optional AllowList — are collected under a catch-all
// group (default "OTHER").
//
// Example
//
//	env := map[string]string{
//		"DB_HOST":  "localhost",
//		"DB_PORT":  "5432",
//		"APP_NAME": "myapp",
//	}
//	groups := grouper.Group(env, grouper.DefaultOptions())
//	// groups["DB"]  → {"DB_HOST": "localhost", "DB_PORT": "5432"}
//	// groups["APP"] → {"APP_NAME": "myapp"}
package grouper
