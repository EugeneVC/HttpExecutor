package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("ADDR")
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatal("Wrong timeout", os.Getenv("TIMEOUT"))
	}

	hybridStorage := NewHybridRepository()
	requestHandler := NewRequestHandler(hybridStorage, 10*time.Duration(timeout)*time.Second)
	server := NewHttpServer(addr, requestHandler)

	server.Start()

	<-stop

	log.Println("\nShutting down the server...")

	server.Stop()

	log.Println("Server gracefully stopped")
}
