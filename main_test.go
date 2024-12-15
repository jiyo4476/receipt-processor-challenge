package main

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestGetPortValid(t *testing.T) {
	t.Setenv("PORT", "80")
	port := getPort()
	assert.Equal(t, "80", port)
}

func TestGetPortInvalid_Letter(t *testing.T) {
	t.Setenv("PORT", "80a")
	port := getPort()
	assert.Equal(t, "8080", port)
}

func TestGetPortInvalid_Negitive(t *testing.T) {
	t.Setenv("PORT", "-1")
	port := getPort()
	assert.Equal(t, "8080", port)
}

func TestGetPortInvalid_OutOfBounds(t *testing.T) {
	t.Setenv("PORT", "65536")
	port := getPort()
	assert.Equal(t, "8080", port)
}
