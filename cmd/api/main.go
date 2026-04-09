package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/is-allowed", isAllowed)

	fmt.Println("Starting server on port :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
