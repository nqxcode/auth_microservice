package async

import (
	"context"
	"log"
)

// Handler handler function type
type Handler func(ctx context.Context) error

// Runner runner interface
type Runner interface {
	Run(ctx context.Context, handler Handler)
}

type runner struct {
}

// NewRunner new runner
func NewRunner() Runner {
	return &runner{}
}

func (r *runner) Run(ctx context.Context, handler Handler) {
	go func() {
		err := handler(ctx)
		if err != nil {
			log.Println("Error in handler:", err)
		}
	}()
}
