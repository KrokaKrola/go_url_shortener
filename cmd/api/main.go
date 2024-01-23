package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/krokakrola/url_shortener/internal/server"
	"github.com/krokakrola/url_shortener/internal/store"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := store.NewDatabase()

	if err != nil {
		log.Fatal("Error initializing database", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal("Error migrating database ", err)
	}

	defer db.Connection.Close(context.Background())

	redis := store.NewRedis()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	s := server.NewServer(db.Connection, redis.Rdb)

	go func() {
		err = s.Start()

		if err != nil {
			panic(err)
		}
	}()

	// Block until a signal is received
	<-stopChan

	log.Println("Shutting down the server...")

	// Attempt to gracefully shut down the server
	if err := s.App.Shutdown(); err != nil {
		log.Println("Error during server shutdown:", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
