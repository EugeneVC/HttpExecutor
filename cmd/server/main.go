package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":8887"
	}

	mapStorage := NewMapStorageRepository()

	s := NewServer(&mapStorage)

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
