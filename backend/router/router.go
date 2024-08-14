package router

import (
	"net/http"

	"github.com/assaidy/personal-blog/backend/handlers"
	"github.com/gorilla/mux"
)

func GetRouter() http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/articles",
		handlers.MakeHttpHandlerFunc(handlers.HandleGetAllArticles)).Methods("GET", "OPTIONS")
	router.HandleFunc("/articles",
		handlers.MakeHttpHandlerFunc(handlers.HandleCreateArticle)).Methods("POST", "OPTIONS")
	router.HandleFunc("/articles",
		handlers.MakeHttpHandlerFunc(handlers.HandleDeleteAllArticles)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/articles/{id}",
		handlers.MakeHttpHandlerFunc(handlers.HandleGetArticleById)).Methods("GET", "OPTIONS")
	router.HandleFunc("/articles/{id}",
		handlers.MakeHttpHandlerFunc(handlers.HandleUpdateArticleById)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/articles/{id}",
		handlers.MakeHttpHandlerFunc(handlers.HandleDeleteArticleById)).Methods("DELETE", "OPTIONS")

	return router
}
