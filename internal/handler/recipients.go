package handler

import (
	"database/sql"
	"net/http"
	"schoolbooks/internal/config"
	"schoolbooks/internal/model"
	"schoolbooks/internal/page"
	"schoolbooks/internal/templates"

	"github.com/go-chi/chi/v5"
)

type RecipientHandler struct {
	DB     *sql.DB
	Config *config.Config
}

func (h *RecipientHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	pd := page.New("Odbiorcy", r, w, h.Config)
	recipients, err := model.ListRecipients(h.DB)
	if err != nil {
		http.Error(w, "Nie udało się załadować odbiorców", http.StatusInternalServerError)
		return
	}
	templates.Recipients(pd, recipients).Render(r.Context(), w)
}

func (h *RecipientHandler) NewRecipientForm(w http.ResponseWriter, r *http.Request) {
	templates.NewRecipientForm().Render(r.Context(), w)
}

func (h *RecipientHandler) EditRecipientForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	rcp, err := model.GetRecipientByID(h.DB, id)
	if err != nil {
		http.Error(w, "Nie udało się załadować odbiorcy", http.StatusInternalServerError)
		return
	}
	templates.EditRecipientForm(rcp).Render(r.Context(), w)
}

func (h *RecipientHandler) Create(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	account := r.FormValue("account")
	description := r.FormValue("description")

	if err := model.CreateRecipient(h.DB, title, account, description); err != nil {
		http.Error(w, "Nie udało się utworzyć odbiorcy", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func (h *RecipientHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	title := r.FormValue("title")
	account := r.FormValue("account")
	description := r.FormValue("description")

	if err := model.UpdateRecipient(h.DB, id, title, account, description); err != nil {
		http.Error(w, "Nie udało się zaktualizować odbiorcy", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func (h *RecipientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := model.DeleteRecipient(h.DB, id); err != nil {
		http.Error(w, "Nie udało się usunąć odbiorcy", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
