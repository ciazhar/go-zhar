package bootstrap

import "context"

type Service interface {
	Start() error
	Shutdown(ctx context.Context) error
	Name() string
}
