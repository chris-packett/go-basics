package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleHeaders(w http.ResponseWriter, r *http.Request) {
	for name, headers := range r.Header {
		for _, header := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, header)
		}
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/headers", handleHeaders)
	log.Fatal(http.ListenAndServe(":8080", limit(mux)))
}
