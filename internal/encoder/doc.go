// Package encoder serializes env maps into different target formats.
//
// Supported formats:
//
//   - shell  — exports suitable for sourcing in bash/zsh (export KEY="value")
//   - docker — plain KEY=value pairs for use with docker --env-file
//   - yaml   — simple YAML mapping (KEY: value)
//
// Example usage:
//
//	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
//	out, err := encoder.Encode(env, encoder.DefaultOptions())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Print(out)
package encoder
