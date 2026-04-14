package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mayukh551/cloudbox/db"
	"github.com/mayukh551/cloudbox/routers"
	"github.com/rs/cors"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = db.Init()

	if err != nil {
		panic(err)
	}

	// utils.Configure("go_logs.log")

	r := routers.Router()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	PORT := os.Getenv("PORT")

	addr := fmt.Sprintf(":%s", PORT)

	fmt.Println("Server listening on port", PORT)
	srv := &http.Server{Addr: addr, Handler: handler}

	// Graceful shutdown of server

	// channel to store OS signals
	quit := make(chan os.Signal, 1)

	// Notify on SIGINT (Ctrl+C) and SIGTERM (sent by Kubernetes on pod termination)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	// block until a signal is received
	sig := <-quit
	log.Printf("Received signal: %s. Initiating graceful shutdown...", sig)

	ctxt, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxt); err != nil {
		log.Fatalf("Server shutdown error: %s\n", err)
	}

	log.Println("Server gracefully stopped")

}
