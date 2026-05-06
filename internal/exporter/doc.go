// Package exporter provides functionality for generating env file templates
// from a differ.Result. It supports multiple output formats:
//
//   - dotenv (.env style KEY=VALUE)
//   - json   (a flat JSON object)
//   - markdown (a two-column table)
//
// Values can optionally be included in the output; by default they are
// redacted so the template is safe to commit to version control.
//
// Example usage:
//
//	opts := exporter.Options{
//		Format:      exporter.FormatDotEnv,
//		IncludeVals: false,
//	}
//	if err := exporter.Export(result, leftEnv, ".env.template", opts); err != nil {
//		log.Fatal(err)
//	}
package exporter
