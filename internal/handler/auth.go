package handler

import (
	"database/sql"
	"net/http"
	"schoolbooks/internal/auth"
	"schoolbooks/internal/config"
	"schoolbooks/internal/model"
	"schoolbooks/internal/page"
	"schoolbooks/internal/session"
	"schoolbooks/internal/templates"
)

type AuthHandler struct {
	DB     *sql.DB
	Config *config.Config
}

func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	pd := page.New("Logowanie", r, w, h.Config)

	templates.Login(pd).Render(r.Context(), w)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	hash, err := model.GetUserPassword(h.DB, email)
	if err != nil {
		session.SetFlash(w, r, "Nieprawidłowy email lub hasło", "error")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !auth.CheckPassword(password, hash) {
		session.SetFlash(w, r, "Nieprawidłowy email lub hasło", "error")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := model.GetUserByEmail(h.DB, email)
	if err != nil {
		session.SetFlash(w, r, "Nieprawidłowy email lub hasło", "error")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	s, _ := session.GetSession(r)
	s.Values["user_id"] = user.ID
	s.Values["name"] = user.Name
	s.Values["role"] = user.Role
	if err := s.Save(r, w); err != nil {
		http.Error(w, "Błąd sesji", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	s, _ := session.GetSession(r)
	s.Options.MaxAge = -1
	s.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
