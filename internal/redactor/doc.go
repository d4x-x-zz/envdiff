// Package redactor masks sensitive values in .env maps before they are
// printed, exported, or compared.
//
// # Usage
//
// Use DefaultOptions to get a pre-configured set of patterns that cover the
// most common secret key names (password, token, secret, key, auth, …):
//
//	masked := redactor.Redact(env, redactor.DefaultOptions())
//
// You can also supply an explicit list of keys to redact regardless of their
// name, or override the mask string:
//
//	opts := redactor.DefaultOptions()
//	opts.ExplicitKeys = []string{"MY_INTERNAL_ID"}
//	opts.Mask = "<hidden>"
//	masked := redactor.Redact(env, opts)
//
// Redact never modifies the original map; it always returns a new one.
package redactor
