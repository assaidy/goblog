package router

import (
	"net/http"

	"github.com/assaidy/goblog/handlers"
	"github.com/assaidy/goblog/repo"
	"github.com/assaidy/goblog/utils"
	"github.com/gorilla/mux"
)

func NewRouter(store repo.Storer) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	userHandler := handlers.NewUserHandler(store)

	router.HandleFunc("/api/register",
		utils.MakeHandlerFunc(userHandler.HandleRegisterUser)).Methods("POST")

	router.HandleFunc("/api/login",
		utils.MakeHandlerFunc(userHandler.HandleLoginUser)).Methods("POST")

	router.HandleFunc("/api/users",
		utils.MakeHandlerFunc(userHandler.HandleGetAllUsers)).Methods("GET")

	router.HandleFunc("/api/users/{id:[0-9]+}",
		utils.MakeHandlerFunc(userHandler.HandleGetUserById)).Methods("GET")

	router.HandleFunc("/api/users/{username}",
		utils.MakeHandlerFunc(userHandler.HandleGetUserByUsername)).Methods("GET")

	router.HandleFunc("/api/users/{id:[0-9]+}",
		utils.MakeHandlerFunc(userHandler.HandleUpdateUserById)).Methods("PUT")

	router.HandleFunc("/api/users/{id:[0-9]+}",
		utils.MakeHandlerFunc(userHandler.HandleDeleteUserById)).Methods("DELETE")

    // other handlers with their routes

	return router
}
