package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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
	timeout,err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err!=nil{
		log.Fatal("Wrong timeout",os.Getenv("TIMEOUT"))
	}

	mapStorage := NewHybridRepository()
	s := NewRequestHandler(&mapStorage,time.Duration(timeout)*time.Second)

	h := &http.Server{Addr: addr, Handler: s}

	go func() {
		log.Printf("Listening on %s\n", addr)

		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Println("\nShutting down the server...")

	h.Shutdown(context.Background())

	log.Println("Server gracefully stopped")
}
