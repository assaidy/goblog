package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/assaidy/personal-blog/backend/helpers"
	"github.com/assaidy/personal-blog/backend/models"
	"github.com/assaidy/personal-blog/backend/repo"
	"github.com/gorilla/mux"
)

var db models.Storager

func init() {
	var err error
	db, err = repo.NewPostgresStore()

	if err != nil {
		log.Fatal("error initializing postgres:", err)
	}
}

func HandleGetAllArticles(w http.ResponseWriter, r *http.Request) error {
	articles, err := db.GetAllArticles()
	if err != nil {
		return err
	}

	return writeJSON(w, articles, http.StatusOK)
}

func HandleDeleteAllArticles(w http.ResponseWriter, r *http.Request) error {
	return db.DeleteAllArticles()
}

func HandleGetArticleById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.NewBadRequestError(fmt.Errorf("invalid id. id must be an integer"))
	}

	article, err := db.GetArticle(id)
	if err != nil {
		return err
	}

	return writeJSON(w, article, http.StatusOK)
}

func HandleCreateArticle(w http.ResponseWriter, r *http.Request) error {
	var articleRequest models.ArticleRequest
	err := json.NewDecoder(r.Body).Decode(&articleRequest)
	if err != nil {
		return models.NewBadRequestError(fmt.Errorf("invalid json request format"))
	}
	defer r.Body.Close()

	article := models.Article{
		Title:       articleRequest.Title,
		Content:     articleRequest.Content,
		PublishDate: time.Now().UTC(),
	}

	if err = db.CreateArticle(&article); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

func HandleDeleteArticleById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.NewBadRequestError(fmt.Errorf("invalid id. id must be an integer"))
	}

	err = db.DeleteArticle(id)
	if err != nil {
		return err
	}

	return nil
}

func HandleUpdateArticleById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.NewBadRequestError(fmt.Errorf("invalid id. id must be an integer"))
	}

	var articleRequest models.ArticleRequest
	err = json.NewDecoder(r.Body).Decode(&articleRequest)
	if err != nil {
		return models.NewBadRequestError(fmt.Errorf("invalid json request format"))
	}
	defer r.Body.Close()

	err = db.UpdateArticle(id, articleRequest.Title, articleRequest.Content)
	if err != nil {
		return err
	}

	return nil
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

func MakeHttpHandlerFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.LogNewRequest(r)

		err := f(w, r)
		if err != nil {
			helpers.LogError(err.Error())

			errorResponse := map[string]string{"error": err.Error()}

			switch err.(type) {
			case *models.DBError:
				errorResponse = map[string]string{"error": "something went wrong"}
				writeJSON(w, errorResponse, http.StatusInternalServerError)
			case *models.NotFoundError:
				writeJSON(w, errorResponse, http.StatusNotFound)
			case *models.BadRequest:
				writeJSON(w, errorResponse, http.StatusBadRequest)
			default:
				errorResponse = map[string]string{"error": "bad request"}
				writeJSON(w, errorResponse, http.StatusBadRequest)
			}
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return models.NewJSONEncodeError(err)
	}
	return nil
}
