package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupLogger(t *testing.T) {
	test_logger := setUpLogger()
	assert.NotNil(t, test_logger, "Logger should not be nil")
}
