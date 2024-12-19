package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	test_logger := getLogger()
	assert.NotNil(t, test_logger, "Logger should not be nil")
}

func TestGetEnv(t *testing.T) {
	t.Setenv("RECEIPT_PROCESSOR_HOSTNAME", "localhost")
	t.Setenv("RECEIPT_PROCESSOR_PORT", "8080")
	test_env := getEnv()
	assert.Equal(t, test_env.HOSTNAME, "localhost", "hostname should be localhost")
	assert.Equal(t, test_env.PORT, "8080", "port should be 8080")
}

func TestGetServer(t *testing.T) {
	test_server := getServer()
	assert.NotNil(t, test_server, "Server should not be nil")
}
