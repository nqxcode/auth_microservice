package tests

import (
	"os"
	"testing"

	"github.com/nqxcode/auth_microservice/internal/logger"
	"github.com/nqxcode/auth_microservice/internal/tracing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func setup() {
	logger.InitNoop()
	tracing.InitNoop()
}

func teardown() {
	// do nothing
}
