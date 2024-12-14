package server

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isPortInUse(port int) bool {
	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return true // Port is in use
	}
	defer ln.Close()
	return false // Port is not in use
}

func TestStartServer(t *testing.T) {
	// Assuming ENV PORT is not set
	Start()
	value, exists := os.LookupEnv("PORT")
	if !exists {
		value = "8080"
	}
	port, err := strconv.Atoi(value)
	if err != nil {
		t.Errorf("Error converting port to integer: %v", err)
	}

	assert.True(t, isPortInUse(port), fmt.Sprintf("Port %d is not in use.", port))
}
