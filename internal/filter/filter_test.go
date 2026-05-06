package filter_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/filter"
)

func makeResult() differ.Result {
	return differ.Result{
		MissingInRight: []differ.Entry{{Key: "DB_HOST"}, {Key: "APP_PORT"}},
		MissingInLeft:  []differ.Entry{{Key: "REDIS_URL"}},
		Mismatched:     []differ.Entry{{Key: "DB_PASS", LeftVal: "secret", RightVal: "other"}},
	}
}

func TestApply_NoOptions(t *testing.T) {
	res := filter.Apply(makeResult(), filter.Options{})
	if len(res.MissingInRight) != 2 {
		t.Errorf("expected 2 MissingInRight, got %d", len(res.MissingInRight))
	}
	if len(res.MissingInLeft) != 1 {
		t.Errorf("expected 1 MissingInLeft, got %d", len(res.MissingInLeft))
	}
	if len(res.Mismatched) != 1 {
		t.Errorf("expected 1 Mismatched, got %d", len(res.Mismatched))
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	res := filter.Apply(makeResult(), filter.Options{OnlyMissing: true})
	if len(res.MissingInRight) != 2 {
		t.Errorf("expected 2 MissingInRight, got %d", len(res.MissingInRight))
	}
	if len(res.Mismatched) != 0 {
		t.Errorf("expected 0 Mismatched, got %d", len(res.Mismatched))
	}
}

func TestApply_OnlyMismatched(t *testing.T) {
	res := filter.Apply(makeResult(), filter.Options{OnlyMismatched: true})
	if len(res.Mismatched) != 1 {
		t.Errorf("expected 1 Mismatched, got %d", len(res.Mismatched))
	}
	if len(res.MissingInRight) != 0 {
		t.Errorf("expected 0 MissingInRight, got %d", len(res.MissingInRight))
	}
}

func TestApply_KeyPrefix(t *testing.T) {
	res := filter.Apply(makeResult(), filter.Options{KeyPrefix: "DB_"})
	if len(res.MissingInRight) != 1 || res.MissingInRight[0].Key != "DB_HOST" {
		t.Errorf("expected only DB_HOST in MissingInRight, got %v", res.MissingInRight)
	}
	if len(res.Mismatched) != 1 || res.Mismatched[0].Key != "DB_PASS" {
		t.Errorf("expected only DB_PASS in Mismatched, got %v", res.Mismatched)
	}
	if len(res.MissingInLeft) != 0 {
		t.Errorf("expected 0 MissingInLeft, got %d", len(res.MissingInLeft))
	}
}

func TestApply_KeyPrefix_NoMatch(t *testing.T) {
	res := filter.Apply(makeResult(), filter.Options{KeyPrefix: "NOPE_"})
	if !res.Clean() {
		t.Error("expected clean result when no keys match prefix")
	}
}
