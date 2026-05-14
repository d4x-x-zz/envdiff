// Package deduper provides utilities for detecting and removing duplicate
// key-value entries when combining multiple .env sources.
//
// A key is considered a duplicate when it has already been encountered in an
// earlier source map.  Whether two entries with the same key but different
// values are treated as duplicates is controlled by Options.SkipValueCheck.
//
// Basic usage:
//
//	import "github.com/yourorg/envdiff/internal/deduper"
//
//	result := deduper.Dedupe(
//		[]map[string]string{base, override},
//		deduper.DefaultOptions(),
//	)
//	fmt.Println(result.Env)     // merged, deduplicated map
//	fmt.Println(result.Removed) // entries that were stripped
package deduper
