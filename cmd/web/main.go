package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux() // servemux = router
	mux.HandleFunc("/", home) // "/" is a catch-all with go's servemux. Subtree path, so matches "/**"
	mux.HandleFunc("/snippet/view", snippetView) // Fixed path - only matched exactly
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on :8080")

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
