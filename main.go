package main

import (
	"fmt"
	"net/http"

	"github.com/mayukh551/cloudbox/db"
	"github.com/mayukh551/cloudbox/routers"
)

const PORT = 7500

func main() {

	err := db.Init()

	if err != nil {
		panic(err)
	}

	// utils.Configure("go_logs.log")

	r := routers.Router()

	addr := fmt.Sprintf(":%d", PORT)

	fmt.Println("Server listening on port", PORT)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		fmt.Println("Error listening on port", PORT)
		panic(err)
	}
}
