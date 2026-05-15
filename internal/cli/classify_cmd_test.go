package cli_test

import (
	"os"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/cli"
)

func writeClassifyTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRunClassify_TextOutput(t *testing.T) {
	path := writeClassifyTempEnv(t, "DB_HOST=localhost\nJWT_SECRET=abc\nAPP_NAME=myapp\n")
	out, err := captureOutput(func() error {
		return cli.RunClassify([]string{path}, nil)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "[database]") {
		t.Error("expected [database] section")
	}
	if !strings.Contains(out, "[auth]") {
		t.Error("expected [auth] section")
	}
}

func TestRunClassify_JSONOutput(t *testing.T) {
	path := writeClassifyTempEnv(t, "DB_HOST=localhost\nHTTP_PORT=8080\n")
	out, err := captureOutput(func() error {
		return cli.RunClassify([]string{"--json", path}, nil)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"database"`) {
		t.Error("expected database key in JSON")
	}
	if !strings.Contains(out, `"network"`) {
		t.Error("expected network key in JSON")
	}
}

func TestRunClassify_MissingFile(t *testing.T) {
	err := cli.RunClassify([]string{"/no/such/file.env"}, nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunClassify_NoArgs(t *testing.T) {
	err := cli.RunClassify([]string{}, nil)
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestRunClassify_OtherCategory(t *testing.T) {
	path := writeClassifyTempEnv(t, "APP_NAME=envdiff\nVERSION=1.0\n")
	out, err := captureOutput(func() error {
		return cli.RunClassify([]string{path}, nil)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "[other]") {
		t.Error("expected [other] category for unmatched keys")
	}
}
