package main

import (
	"context"

	"github.com/yebrai/go-tasks-microservice/internal/task/bootstrap"
	"github.com/yebrai/go-tasks-microservice/pkg/runner"
)

func main() {
	runner.Run(func(ctx context.Context) runner.Runner {
		return bootstrap.NewService()
	})
}
