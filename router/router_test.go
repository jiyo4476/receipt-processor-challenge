package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	test_router := SetUpRouter()
	assert.NotNil(t, test_router, "Router should not be nil")
}

func TestSetupLogger(t *testing.T) {
	test_logger := SetUpRouter()
	assert.NotNil(t, test_logger, "Logger should not be nil")
}
