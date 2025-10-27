package args

import (
	"flag"
	"fmt"
	"os"
)

// Parses args from [os.Args] with [flag].
func ParseArgs() (*Args, error) {
	return Parse(os.Args)
}

// Parses args with [flag].
func Parse(arguments []string) (*Args, error) {
	flagSet := flag.NewFlagSet(arguments[0], flag.ContinueOnError)
	configPath := flagSet.String("config", "", "Configuration file url")
	flagSet.Parse(arguments[1:])
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
