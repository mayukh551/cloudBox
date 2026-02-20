package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	err = http.ListenAndServe(addr, handler)
	if err != nil {
		fmt.Println("Error listening on port", PORT)
		panic(err)
	}
}
