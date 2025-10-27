package args

import (
	"fmt"
	"os"
)

// Scans the args for invalid settings.
func (a *Args) validateConfig() error {
	fileInfo, err := os.Stat(a.Config)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("config should be a file")
	}
	return nil
}
