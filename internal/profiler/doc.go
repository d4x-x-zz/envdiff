// Package profiler provides statistical analysis of env maps.
//
// It computes key counts, value length statistics, type distribution
// (bool, int, float, url, string), empty-value counts, and density
// (fraction of keys with non-empty values).
//
// Example:
//
//	env := map[string]string{"PORT": "8080", "DEBUG": "true", "SECRET": ""}
//	p := profiler.Analyze(env, profiler.DefaultOptions())
//	fmt.Printf("density: %.2f\n", p.Density)
//	fmt.Printf("types: %v\n", profiler.SortedTypes(p))
package profiler
