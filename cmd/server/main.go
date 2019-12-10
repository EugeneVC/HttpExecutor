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
	timeoutSecond, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatal("Wrong timeout", os.Getenv("TIMEOUT"))
	}

	hybridTaskRepository := NewHybridTaskRepository()
	requestHandler := NewRequestHandler(hybridTaskRepository, time.Duration(timeoutSecond)*time.Second)
	server := NewHttpServer(addr, requestHandler)

	server.Start()

	<-stop

	log.Println("\nShutting down the server...")

	server.Stop()

	log.Println("Server gracefully stopped")
}
