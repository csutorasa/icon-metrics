package publisher

import (
	"context"
	"io"
)

type Publisher interface {
	io.Closer
	Start() error
	Stop(context context.Context) error
}
