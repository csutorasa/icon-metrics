package args

import (
	"flag"
	"fmt"
)

// Parses args from command line options with flag
func ParseArgs(arguments []string) (*Args, error) {
	configPath := flag.String("config", "", "Configuration file url")
	flag.CommandLine.Parse(arguments)
	args := &Args{
		Config: *configPath,
	}
	args.setDefaults()
	err := args.validateConfig()
	if err != nil {
		return args, fmt.Errorf("invalid args: %w", err)
	}
	return args, nil
}
