package support

import (
	"context"
	"log"

	"github.com/nqxcode/auth_microservice/internal/service/async"
)

// AsyncRunnerFake async runner fake type
type AsyncRunnerFake struct{}

// NewAsyncRunnerFake new fake async runner
func NewAsyncRunnerFake() async.Runner {
	return &AsyncRunnerFake{}
}

// Run run handler
func (runner *AsyncRunnerFake) Run(ctx context.Context, handler async.Handler) {
	err := handler(ctx)
	if err != nil {
		log.Println("Error in handler:", err)
	}
}
