package differ

// Kind describes the type of difference found between two env maps.
type Kind string

const (
	MissingInRight Kind = "missing_right"
	MissingInLeft  Kind = "missing_left"
	ValueMismatch  Kind = "value_mismatch"
)

// Result holds a single difference between two env maps.
type Result struct {
	Key      string
	Kind     Kind
	LeftVal  string
	RightVal string
}

// Diff compares two env maps and returns all differences.
func Diff(left, right map[string]string) []Result {
	var results []Result

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			results = append(results, Result{
				Key:     k,
				Kind:    MissingInRight,
				LeftVal: lv,
			})
			continue
		}
		if lv != rv {
			results = append(results, Result{
				Key:      k,
				Kind:     ValueMismatch,
				LeftVal:  lv,
				RightVal: rv,
			})
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			results = append(results, Result{
				Key:      k,
				Kind:     MissingInLeft,
				RightVal: rv,
			})
		}
	}

	return results
}
