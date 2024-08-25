package async

import (
	"context"
	"fmt"
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
		defer func() {
			if rr := recover(); rr != nil {
				err, ok := rr.(error)
				if !ok {
					err = fmt.Errorf("%v", rr)
				}
				log.Println("Panic in handler:", err)
			}
		}()

		err := handler(ctx)
		if err != nil {
			log.Println("Error in handler:", err)
		}
	}()
}
