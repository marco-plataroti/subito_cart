package main

import (
	"log"
	"net/http"

	"subito-cart/router" // replace with your module name

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	router.RegisterRoutes(r)

	const addr = ":9090"
	log.Printf("Server is listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
