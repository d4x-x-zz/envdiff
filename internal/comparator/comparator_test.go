package comparator_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
)

func makeFileEnv() comparator.FileEnv {
	return comparator.FileEnv{
		"dev": {"DB_HOST": "localhost", "DB_PORT": "5432", "SECRET": "dev-secret"},
		"staging": {"DB_HOST": "staging.db", "DB_PORT": "5432", "STAGING_ONLY": "yes"},
		"prod": {"DB_HOST": "prod.db", "DB_PORT": "5433", "SECRET": "prod-secret"},
	}
}

func TestCompare_FilesAreSorted(t *testing.T) {
	r := comparator.Compare(makeFileEnv())
	if r.Files[0] != "dev" || r.Files[1] != "prod" || r.Files[2] != "staging" {
		t.Fatalf("unexpected file order: %v", r.Files)
	}
}

func TestCompare_KeysAreSorted(t *testing.T) {
	r := comparator.Compare(makeFileEnv())
	if r.Statuses[0].Key != "DB_HOST" {
		t.Fatalf("expected first key DB_HOST, got %s", r.Statuses[0].Key)
	}
}

func TestCompare_UniformKey(t *testing.T) {
	r := comparator.Compare(makeFileEnv())
	var portStatus *comparator.KeyStatus
	for i := range r.Statuses {
		if r.Statuses[i].Key == "DB_PORT" {
			portStatus = &r.Statuses[i]
			break
		}
	}
	if portStatus == nil {
		t.Fatal("DB_PORT not found")
	}
	if !portStatus.Uniform {
		t.Error("expected DB_PORT to be uniform")
	}
	if len(portStatus.Missing) != 0 {
		t.Errorf("expected no missing files for DB_PORT, got %v", portStatus.Missing)
	}
}

func TestCompare_MismatchedKey(t *testing.T) {
	r := comparator.Compare(makeFileEnv())
	for _, s := range r.Statuses {
		if s.Key == "DB_HOST" {
			if s.Uniform {
				t.Error("expected DB_HOST to be non-uniform")
			}
			return
		}
	}
	t.Fatal("DB_HOST not found")
}

func TestCompare_MissingKey(t *testing.T) {
	r := comparator.Compare(makeFileEnv())
	for _, s := range r.Statuses {
		if s.Key == "STAGING_ONLY" {
			if len(s.Missing) != 2 {
				t.Errorf("expected 2 missing files, got %d", len(s.Missing))
			}
			return
		}
	}
	t.Fatal("STAGING_ONLY not found")
}

func TestReport_TotalMissing(t *testing.T) {
	r := comparator.Compare(makeFileEnv())
	// SECRET missing in staging, STAGING_ONLY missing in dev+prod
	if r.TotalMissing() != 2 {
		t.Errorf("expected TotalMissing 2, got %d", r.TotalMissing())
	}
}

func TestReport_TotalMismatched(t *testing.T) {
	r := comparator.Compare(makeFileEnv())
	// DB_HOST and DB_PORT differ; DB_PORT is uniform so only DB_HOST counts
	// SECRET is missing in staging so it doesn't count as mismatched
	if r.TotalMismatched() != 1 {
		t.Errorf("expected TotalMismatched 1, got %d", r.TotalMismatched())
	}
}

func TestCompare_EmptyInput(t *testing.T) {
	r := comparator.Compare(comparator.FileEnv{})
	if len(r.Statuses) != 0 {
		t.Error("expected no statuses for empty input")
	}
}
