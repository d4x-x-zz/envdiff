// Package classifier categorises env keys into semantic buckets
// such as database, auth, network, storage, etc.
package classifier

import "strings"

// Category represents a semantic group for an env key.
type Category string

const (
	CategoryDatabase Category = "database"
	CategoryAuth     Category = "auth"
	CategoryNetwork  Category = "network"
	CategoryStorage  Category = "storage"
	CategoryLogging  Category = "logging"
	CategoryOther    Category = "other"
)

// Result holds the classification output.
type Result struct {
	// Categories maps each category to the list of keys that belong to it.
	Categories map[Category][]string
}

// Options controls classifier behaviour.
type Options struct {
	// IncludeOther, when true, includes keys that don't match any known
	// category under the "other" bucket.
	IncludeOther bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{IncludeOther: true}
}

var categoryPatterns = map[Category][]string{
	CategoryDatabase: {"DB_", "DATABASE_", "POSTGRES", "MYSQL", "MONGO", "REDIS", "DSN"},
	CategoryAuth:     {"AUTH_", "JWT_", "SECRET", "TOKEN", "PASSWORD", "PASSWD", "API_KEY"},
	CategoryNetwork:  {"HOST", "PORT", "URL", "ADDR", "ENDPOINT", "PROXY", "TLS_", "SSL_"},
	CategoryStorage:  {"S3_", "BUCKET", "STORAGE_", "VOLUME", "DISK_", "PATH", "DIR"},
	CategoryLogging:  {"LOG_", "LOGGER_", "DEBUG", "VERBOSE", "TRACE", "SENTRY_"},
}

// Classify assigns each key in env to a semantic category.
func Classify(env map[string]string, opts Options) Result {
	result := Result{Categories: make(map[Category][]string)}

	for key := range env {
		upper := strings.ToUpper(key)
		matched := false

		for cat, patterns := range categoryPatterns {
			for _, p := range patterns {
				if strings.Contains(upper, p) {
					result.Categories[cat] = append(result.Categories[cat], key)
					matched = true
					break
				}
			}
			if matched {
				break
			}
		}

		if !matched && opts.IncludeOther {
			result.Categories[CategoryOther] = append(result.Categories[CategoryOther], key)
		}
	}

	return result
}
