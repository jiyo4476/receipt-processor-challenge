package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/jiyo4476/receipt-processor-challenge/router"
)

func Start() {
	cur_router := router.SetUpRouter()

	value, exists := os.LookupEnv("PORT")
	if !exists {
		value = "8080"
	}

	server := &http.Server{
		Addr:    ":" + value,
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
