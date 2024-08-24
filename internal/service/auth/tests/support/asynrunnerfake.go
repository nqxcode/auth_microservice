package support

import (
	"context"
	"github.com/nqxcode/auth_microservice/internal/service/async"
	"log"
)

type AsyncRunnerFake struct{}

func NewAsyncRunnerFake() async.Runner {
	return &AsyncRunnerFake{}
}

func (runner *AsyncRunnerFake) Run(ctx context.Context, handler async.Handler) {
	err := handler(ctx)
	if err != nil {
		log.Println("Error in handler:", err)
	}
}
