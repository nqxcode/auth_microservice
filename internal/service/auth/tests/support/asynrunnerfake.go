package support

import "github.com/nqxcode/auth_microservice/internal/service/async"

type AsyncRunnerFake struct{}

func NewAsyncRunnerFake() async.Runner {
	return &AsyncRunnerFake{}
}

func (runner *AsyncRunnerFake) Run(handler func()) {
	handler()
}
