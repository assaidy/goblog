package main

import (
	"log"
	"net/http"

	"github.com/assaidy/personal-blog/backend/router"
)

const Port = ":6868"

func main() {
	router := router.GetRouter()

	log.Printf("Starting server at port %s...\n", Port)

	log.Fatal(http.ListenAndServe(Port, router))
}
