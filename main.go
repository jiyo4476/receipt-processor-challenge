package main

import (
	"log"

	"github.com/jiyo4476/receipt-processor-challenge/server"
	"github.com/jiyo4476/receipt-processor-challenge/spec"
)

func main() {
	// Load specs in globally accessible variable
	if err := spec.PrintSpec("api.yml"); err != nil {
		log.Printf("Error loading spec: %v", err) // Log the error
		return
	}

	server.Start()
}
