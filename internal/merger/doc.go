// Package merger combines multiple parsed .env maps into a single unified map.
//
// It supports three conflict-resolution strategies:
//
//   - StrategyFirst: the first file to define a key wins (default).
//   - StrategyLast:  the last file to define a key wins.
//   - StrategyError: returns an error when the same key has differing values
//     across files.
//
// Typical usage:
//
//	result, err := merger.Merge([]map[string]string{base, override}, merger.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
package merger
