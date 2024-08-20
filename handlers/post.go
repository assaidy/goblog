package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/assaidy/goblog/models"
	"github.com/assaidy/goblog/repo"
	"github.com/assaidy/goblog/utils"
	"github.com/gorilla/mux"
)

type PostHandler struct {
	store repo.Storer
}

func NewPostHandler(store repo.Storer) *PostHandler {
	return &PostHandler{store: store}
}

func (h *PostHandler) HandleGetAllPosts(w http.ResponseWriter, r *http.Request) error {
	posts, err := h.store.GetAllPosts()
	if err != nil {
		return err
	}
	return utils.WriteJSON(w, http.StatusOK, posts)
}

func (h *PostHandler) HandleGetAllPostsByUser(w http.ResponseWriter, r *http.Request) error {
	userId, _ := strconv.Atoi(mux.Vars(r)["userId"])
	posts, err := h.store.GetAllPostsByAuthor(userId)
	if err != nil {
		return err
	}
	return utils.WriteJSON(w, http.StatusOK, posts)
}

func (h *PostHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) error {
	var postReq models.PostCreateOrUpdateRequest

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&postReq); err != nil {
		r.Body.Close()
		return utils.InvalidJSON()
	}
	defer r.Body.Close()

	// Trim input fields
	postReq.Title = strings.TrimSpace(postReq.Title)
	postReq.Content = strings.TrimSpace(postReq.Content)

	// Check for required fields
	if postReq.Title == "" || postReq.Content == "" {
		return utils.InvalidRequestData([]string{"title and content are required"})
	}

	// Retrieve userId from context
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		return utils.UnAuthorized(fmt.Errorf("user ID missing or invalid"))
	}

	// Check if the author ID matches the user ID
	if userId != postReq.AuthorId {
		return utils.UnAuthorized(fmt.Errorf("your user ID %d does not match the author ID %d", userId, postReq.AuthorId))
	}

	// Create the post
	post := models.Post{
		Title:     postReq.Title,
		Content:   postReq.Content,
		AuthorId:  postReq.AuthorId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Store the post
	postResp, err := h.store.CreatePost(&post)
	if err != nil {
		return err
	}

	// Respond with the created post
	return utils.WriteJSON(w, http.StatusCreated, postResp)
}

func (h *PostHandler) HandleGetPostById(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	post, err := h.store.GetPostById(id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, post)
}

func (h *PostHandler) HandleUpdatePostById(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var updateReq models.PostCreateOrUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		return utils.InvalidJSON()
	}
	defer r.Body.Close()

	// Retrieve userId from context
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		return utils.UnAuthorized(fmt.Errorf("user ID missing or invalid"))
	}

	// Check if the author ID matches the user ID
	if userId != updateReq.AuthorId {
		return utils.UnAuthorized(fmt.Errorf("your user ID %d does not match the author ID %d", userId, updateReq.AuthorId))
	}

	post, err := h.store.UpdatePostById(id, &updateReq)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, post)
}

func (h *PostHandler) HandleDeletePostById(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Retrieve userId from context
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		return utils.UnAuthorized(fmt.Errorf("user ID missing or invalid"))
	}

	// Fetch the post to ensure it exists and to check authorization
	post, err := h.store.GetPostById(id)
	if err != nil {
		return utils.NotFound(fmt.Errorf("post with id %d not found", id))
	}

	// Check if the user is the author of the post
	if post.AuthorId != userId {
		return utils.UnAuthorized(fmt.Errorf("you are not authorized to delete this post"))
	}

	// Proceed with deletion
	if err := h.store.DeletePostById(id, userId); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, nil)
}
