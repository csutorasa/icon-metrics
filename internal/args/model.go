package args

import (
	"fmt"
)

// Command line arguments
type Args struct {
	Config string
}

func (a *Args) String() string {
	return fmt.Sprintf("{config=%s}", a.Config)
}
