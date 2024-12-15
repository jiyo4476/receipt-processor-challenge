package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/jiyo4476/receipt-processor-challenge/router"
	"github.com/jiyo4476/receipt-processor-challenge/spec"
)

func getPort() string {
	value, exists := os.LookupEnv("PORT")
	if !exists {
		fmt.Printf("Environment Variable PORT not defined, using default port 8080")
		return "8080"
	} else {
		port, err := strconv.Atoi(value)
		if err != nil || port < 0 || port > 65535 {
			fmt.Printf("%q is not a valid port number, using default.", value)
			//value = "8080"
			return "8080"
		} else {
			return value
		}
	}
}

func main() {
	// Load specs in globally accessible variable
	if err := spec.PrintSpec("api.yml"); err != nil {
		log.Printf("Error loading spec: %v", err) // Log the error
		return
	}

	cur_router := router.SetUpRouter()

	port := getPort()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: cur_router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("received interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Server Closed:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpectedly")
		}
	}

	log.Println("Server exiting")
}
