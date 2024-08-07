package main

import (
	"crawler-visa/router"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	router.RegisterRouters(r)
	http.Handle("/", r)
	log.Println("Server is starting on localhost:9010...")
	log.Fatal(http.ListenAndServe("localhost:9010", r))
}
