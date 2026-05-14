// Package cloner provides a utility for producing deep copies of env maps
// with optional key and value transformations applied during the copy.
//
// Supported transformations:
//   - KeyPrefix / KeySuffix — wrap every key with a static string
//   - UppercaseKeys         — normalise keys to UPPER_CASE
//   - OmitEmpty             — drop keys with empty values
//   - ValueTransform        — apply an arbitrary func(string) string to values
//
// Example:
//
//	out := cloner.Clone(src, cloner.Options{
//		KeyPrefix:     "STAGING_",
//		UppercaseKeys: true,
//		OmitEmpty:     true,
//	})
package cloner
