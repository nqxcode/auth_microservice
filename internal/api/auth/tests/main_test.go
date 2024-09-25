package tests

import (
	"os"
	"testing"

	"github.com/nqxcode/auth_microservice/internal/logger"
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
