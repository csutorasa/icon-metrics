package args

import (
	"os"
	"path/filepath"
)

func (a *Args) setDefaults() {
	if a.Config == "" {
		dir := filepath.Dir(os.Args[0])
		a.Config = filepath.Join(dir, "config.yml")
	}
}
