package async

type Runner interface {
	Run(handler func())
}

type runner struct {
}

func NewRunner() Runner {
	return &runner{}
}

func (r *runner) Run(handler func()) {
	go handler()
}
