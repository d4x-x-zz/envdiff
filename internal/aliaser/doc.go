// Package aliaser copies env values from source keys to destination alias keys.
//
// This is useful when migrating from one key naming convention to another while
// maintaining backward compatibility, or when a service expects a different key
// name than the one stored in your canonical .env file.
//
// Example usage:
//
//	opts := aliaser.DefaultOptions()
//	opts.Aliases = map[string][]string{
//	    "DB_HOST": {"DATABASE_HOST", "POSTGRES_HOST"},
//	}
//	out, err := aliaser.Alias(env, opts)
package aliaser
