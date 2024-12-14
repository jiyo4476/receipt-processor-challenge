package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	test_router := SetUpRouter()
	assert.NotNil(t, test_router, "Router should not be nil")
}
