package args_test

import (
	"flag"
	"os"
	"testing"

	"github.com/csutorasa/icon-metrics/internal/args"
)

func TestParseArgs(t *testing.T) {
	currentFile := os.Args[0]
	flag.CommandLine = flag.NewFlagSet(currentFile, flag.ExitOnError)
	args, err := args.ParseArgs([]string{
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
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	_, err := args.ParseArgs([]string{
		"-config", "file-that-does-not-exist.yml",
	})
	if err == nil {
		t.Fatalf("error is expected but got none")
	}
}
