package router

import (
	"net/http"

	// "github.com/assaidy/goblog/handlers"
	"github.com/assaidy/goblog/repo"
	"github.com/gorilla/mux"
)

func NewRouter(store repo.Storer) http.Handler {
	router := mux.NewRouter()

	// userHandler := handlers.NewUserHandler(store)

	return router
}
