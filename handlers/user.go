package handlers

import (
	"github.com/assaidy/goblog/repo"
)

type UserHandler struct {
	store repo.Storer
}

func NewUserHandler(store repo.Storer) *UserHandler {
	return &UserHandler{store: store}
}

// handle user here 
