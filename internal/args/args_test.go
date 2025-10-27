package args_test

import (
	"os"
	"testing"

	"github.com/csutorasa/icon-metrics/internal/args"
)

func TestParseArgs(t *testing.T) {
	currentFile := os.Args[0]
	args, err := args.Parse([]string{
		currentFile,
		"-config", currentFile,
	})
	if err != nil {
		t.Fatalf("error is not expected but got %s", err.Error())
	}
	if args.Config != currentFile {
		t.Logf("Expected config /etc/icon-metrics/config.yml but got %s", args.Config)
		t.Fail()
	}
}

func TestParseInvalidArgs(t *testing.T) {
	_, err := args.Parse([]string{
		os.Args[0],
		"-config", "file-that-does-not-exist.yml",
	})
	if err == nil {
		t.Fatalf("error is expected but got none")
	}
}
