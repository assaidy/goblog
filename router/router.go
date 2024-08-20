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

	// Create a protected subrouter with /api prefix
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(utils.JWTAuthMiddleware)

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

	protected.HandleFunc("/api/users/{id:[0-9]+}",
		utils.MakeHandlerFunc(userHandler.HandleUpdateUserById)).Methods("PUT")
	protected.HandleFunc("/api/users/{id:[0-9]+}",
		utils.MakeHandlerFunc(userHandler.HandleDeleteUserById)).Methods("DELETE")

	postHandler := handlers.NewPostHandler(store)

	router.HandleFunc("/api/posts",
		utils.MakeHandlerFunc(postHandler.HandleGetAllPosts)).Methods("GET")
    router.HandleFunc("/api/users/{userId:[0-9]+}/posts",
        utils.MakeHandlerFunc(postHandler.HandleGetAllPostsByUser)).Methods("GET")
	router.HandleFunc("/api/posts/{id:[0-9]+}",
		utils.MakeHandlerFunc(postHandler.HandleGetPostById)).Methods("GET")

	protected.HandleFunc("/posts",
		utils.MakeHandlerFunc(postHandler.HandleCreatePost)).Methods("POST")
	protected.HandleFunc("/posts/{id:[0-9]+}",
		utils.MakeHandlerFunc(postHandler.HandleUpdatePostById)).Methods("PUT")
	protected.HandleFunc("/posts/{id:[0-9]+}",
		utils.MakeHandlerFunc(postHandler.HandleDeletePostById)).Methods("DELETE")

	return router
}
