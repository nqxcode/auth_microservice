package async

import (
	"context"
	"log"
)

type Handler func(ctx context.Context) error

type Runner interface {
	Run(ctx context.Context, handler Handler)
}

type runner struct {
}

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
