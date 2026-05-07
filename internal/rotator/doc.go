// Package rotator implements key rotation for .env maps.
//
// Key rotation is the process of renaming environment variable keys,
// typically when migrating to a new naming convention or deprecating
// old key names.
//
// Usage:
//
//	rotations := []rotator.Rotation{
//		{OldKey: "DB_HOST", NewKey: "DATABASE_HOST"},
//	}
//	opts := rotator.DefaultOptions()
//	newEnv, err := rotator.Rotate(env, rotations, opts)
//
// Strategies:
//   - StrategyRemove  – deletes the old key (default)
//   - StrategyDeprecate – blanks the old key value
//   - StrategyKeep    – leaves the old key untouched
package rotator
