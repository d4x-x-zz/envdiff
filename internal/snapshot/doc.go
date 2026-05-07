// Package snapshot provides save/load/compare functionality for .env file snapshots.
//
// A snapshot captures the full key-value state of an env map at a specific
// point in time and persists it as a JSON file. Snapshots can later be loaded
// and compared to detect keys that were added, removed, or changed between
// two points in time — useful for auditing environment drift.
//
// Basic usage:
//
//	// Save current state
//	snapshot.Save(".env.snap", "before-deploy", env)
//
//	// Load a previous snapshot
//	snap, _ := snapshot.Load(".env.snap")
//
//	// Compare old vs new
//	result := snapshot.Compare(snap, newSnap)
//	if !result.Clean() {
//		fmt.Println("drift detected")
//	}
package snapshot
