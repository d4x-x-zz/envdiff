// Package classifier assigns env keys to semantic categories such as
// database, auth, network, storage, and logging.
//
// Usage:
//
//	env := map[string]string{
//		"DB_HOST":    "localhost",
//		"JWT_SECRET": "s3cr3t",
//		"HTTP_PORT":  "8080",
//	}
//
//	res := classifier.Classify(env, classifier.DefaultOptions())
//	for cat, keys := range res.Categories {
//		fmt.Printf("%s: %v\n", cat, keys)
//	}
//
// Categories are matched by checking whether a key (uppercased) contains
// any of the well-known pattern strings associated with that category.
// Keys that don't match any pattern are placed in the "other" bucket
// when Options.IncludeOther is true (the default).
package classifier
