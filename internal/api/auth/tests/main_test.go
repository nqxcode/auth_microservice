package tests

import (
	"github.com/nqxcode/auth_microservice/internal/logger"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func setup() {
	logger.InitNoop()
}

func teardown() {
	// do nothing
}
