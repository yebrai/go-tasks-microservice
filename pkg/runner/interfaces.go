package runner

import "context"

type Config interface {
	Unmarshal(v interface{}) error
}

type Runner interface {
	Name() string
	Run(ctx context.Context, config Config) error
	Cleanup() error
}
