package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/assaidy/goblog/models"
	"github.com/assaidy/goblog/repo"
	"github.com/assaidy/goblog/utils"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	store repo.Storer
}

func NewUserHandler(store repo.Storer) *UserHandler {
	return &UserHandler{store: store}
}

func (h *UserHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.store.GetAllUsers()
	if err != nil {
		return err
	}
	return utils.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) error {
	var registerReq models.UserRegisterOrUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		return utils.InvalidJSON()
	}
	defer r.Body.Close()

	// Trim input fields
	registerReq.Username = strings.TrimSpace(registerReq.Username)
	registerReq.Email = strings.TrimSpace(registerReq.Email)
	registerReq.Password = strings.TrimSpace(registerReq.Password)

	// Check for required fields
	if registerReq.Username == "" || registerReq.Email == "" || registerReq.Password == "" {
		return utils.InvalidRequestData([]string{"Username, email, and password are required"})
	}

	// Validate user input
	if validationErrors, err := utils.ValidateRegisterUser(registerReq.Username, registerReq.Email, h.store); err != nil {
		return err
	} else if len(validationErrors) > 0 {
		return utils.InvalidRequestData(validationErrors)
	}

	user := models.User{
		FullName: registerReq.FullName,
		Username: registerReq.Username,
		Email:    registerReq.Email,
		Password: registerReq.Password,
		Bio:      registerReq.Bio,
		JoinedAt: time.Now().UTC(),
	}

	userResp, err := h.store.CreateUser(&user)
	if err != nil {
		return err
	}

	// Clear password before sending response
	userResp.Password = ""

	return utils.WriteJSON(w, http.StatusCreated, userResp)
}

func (h *UserHandler) HandleLoginUser(w http.ResponseWriter, r *http.Request) error {
	var loginReq models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		log.Println(err)
		return utils.InvalidJSON()
	}
	defer r.Body.Close()

	// Trim input fields
	loginReq.Username = strings.TrimSpace(loginReq.Username)
	loginReq.Password = strings.TrimSpace(loginReq.Password)

	// Check for required fields
	if loginReq.Username == "" || loginReq.Password == "" {
		return utils.InvalidRequestData([]string{"Username and password are required"})
	}

	user, err := utils.AuthenticateUser(loginReq, h.store)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) HandleGetUserById(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := h.store.GetUserById(id)
	if err != nil {
		return err
	}
	user.Password = ""

	return utils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) HandleGetUserByUsername(w http.ResponseWriter, r *http.Request) error {
	username := mux.Vars(r)["username"]

	user, err := h.store.GetUserByUsername(username)
	if err != nil {
		return err
	}
	user.Password = ""

	return utils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) HandleUpdateUserById(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var updateReq models.UserRegisterOrUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		return utils.InvalidJSON()
	}
	defer r.Body.Close()

	// Authorize user
	if err := utils.AuthorizeUser(id, r); err != nil {
		return err
	}

	user, err := h.store.UpdateUserById(id, &updateReq)
	if err != nil {
		return err
	}
	user.Password = ""

	return utils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) HandleDeleteUserById(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Authorize user
	if err := utils.AuthorizeUser(id, r); err != nil {
		return err
	}

	if err := h.store.DeleteUserById(id); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, nil)
}

