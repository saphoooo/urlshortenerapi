package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var dbDriver *redisDriver

type url struct {
	URL string `json:"url,omitempty"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/new", shortenerHandler).Methods("POST")
	r.HandleFunc("/api/v1/{url}", originalHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
