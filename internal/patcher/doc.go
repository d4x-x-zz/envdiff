// Package patcher provides functionality to apply declarative patches
// to an env map produced by the parser.
//
// A patch is a slice of Patch values, each describing a single set or
// delete operation on a named key.  Apply always works on a copy of the
// source map so the original is never mutated.
//
// Example:
//
//	patches := []patcher.Patch{
//		{Op: patcher.OpSet,    Key: "APP_ENV",  Value: "staging"},
//		{Op: patcher.OpDelete, Key: "DEBUG"},
//	}
//	result, err := patcher.Apply(src, patches, patcher.DefaultOptions())
package patcher
