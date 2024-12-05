package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	router := SetUpRouter()
	assert.NotNil(t, router, "Router should not be nil")
}
