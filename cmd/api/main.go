package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url_shortener_2/internal/server"
)

func main() {
	server := server.NewServer()

	// Channel to receive interrupt signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server started on address %s\n", server.Addr)
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Block until a signal is received
	<-stopChan

	fmt.Println("Shutting down the server...")

	// Create a context with a timeout to allow outstanding requests to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Error during server shutdown:", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}
}
