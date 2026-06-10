package handler

import (
	"database/sql"
	"net/http"
	"schoolbooks/internal/auth"
	"schoolbooks/internal/model"
	"schoolbooks/internal/page"
	"schoolbooks/internal/templates"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	pd := page.New("Admin", r, w)
	users, err := model.ListUsers(h.DB)
	if err != nil {
		http.Error(w, "could not load users", http.StatusInternalServerError)
		return
	}
	templates.Users(pd, users).Render(r.Context(), w)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	role := r.FormValue("role")

	hash, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "could not hash password", http.StatusInternalServerError)
		return
	}

	_, err = model.CreateUser(h.DB, name, email, hash, role)
	if err != nil {
		http.Error(w, "could not create user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
