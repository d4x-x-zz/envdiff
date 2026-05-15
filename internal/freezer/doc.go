// Package freezer captures a point-in-time snapshot of an env map and
// produces a stable SHA-256 fingerprint for change detection.
//
// Usage:
//
//	f := freezer.Freeze(env, freezer.DefaultOptions())
//	fmt.Println(f.Fingerprint)
//
//	// later…
//	if f.Changed(updatedEnv, freezer.DefaultOptions()) {
//		fmt.Println("env has changed since freeze")
//	}
//
// Options let you restrict the freeze to a key prefix and choose whether
// values are included in the fingerprint (keys-only mode is useful when
// you only care about structural changes).
package freezer
