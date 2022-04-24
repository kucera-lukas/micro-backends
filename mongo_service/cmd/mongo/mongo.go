package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Panicf("Failed to listen on %d: %v", 8081, err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello, %s from mongo service!", r.URL.Path[1:])
	if err != nil {
		log.Printf("Failed to greet from mongo service: %v", err)
	}
}
