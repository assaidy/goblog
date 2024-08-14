package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/assaidy/personal-blog/backend/log"
	"github.com/assaidy/personal-blog/backend/types"
)

func MakeHttpHandlerFunc(f types.ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.NewRequest(r)

		err := f(w, r)
		if err != nil { // TODO: handle varias error
			log.Error(err.Error(), http.StatusInternalServerError)

			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleGetAllArticles(w http.ResponseWriter, r *http.Request) error {
	// return writeJSON(w, "HandleGetAllArticles")
	return fmt.Errorf("this is an error")
}

func HandleDeleteAllArticles(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, "HandleDeleteAllArticles")
}

func HandleGetArticleById(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, "HandleGetArticleById")
}

func HandleCreateArticle(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, "HandleCreateArticle")
}

func HandleDeleteArticleById(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, "HandleDeleteArticle")
}

func HandleEditArticleById(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, "HandleEditArticle")
}

func writeJSON(w http.ResponseWriter, v any) error {
	return json.NewEncoder(w).Encode(v)
}
