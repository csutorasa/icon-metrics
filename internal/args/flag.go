package args

import (
	"flag"
	"os"
	"path/filepath"
)

// Parses args from command line options with flag
func ParseArgs() *Args {
	configPath := flag.String("config", "", "Configuration file url")
	flag.Parse()
	var config string
	if *configPath == "" {
		dir := filepath.Dir(os.Args[0])
		config = filepath.Join(dir, "config.yml")
	} else {
		config = *configPath
	}
	return &Args{
		Config: config,
	}
}
