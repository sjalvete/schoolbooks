package handler

import (
	"database/sql"
	"net/http"
	"schoolbooks/internal/auth"
	"schoolbooks/internal/config"
	"schoolbooks/internal/model"
	"schoolbooks/internal/page"
	"schoolbooks/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	DB     *sql.DB
	Config *config.Config
}

func (h *PaymentHandler) UserList(w http.ResponseWriter, r *http.Request) {
	pd := page.New("Płatności", r, w, h.Config)

	user, _ := auth.GetUser(r)
	events, err := model.ListEventAttendance(h.DB, user.ID)
	if err != nil {
		http.Error(w, "Nie udało się załadować płatności", http.StatusInternalServerError)
		return
	}

	templates.UserPayments(pd, model.BucketEvents(events)).Render(r.Context(), w)
}

func (h *PaymentHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	pd := page.New("Płatności", r, w, h.Config)

	events, err := model.ListEventAttendance(h.DB, 0)
	if err != nil {
		http.Error(w, "Nie udało się załadować płatności", http.StatusInternalServerError)
		return
	}

	templates.AdminPayments(pd, model.BucketEvents(events)).Render(r.Context(), w)
}

func (h *PaymentHandler) SetAmountPaid(w http.ResponseWriter, r *http.Request) {
	eventID, err1 := strconv.Atoi(chi.URLParam(r, "eventID"))
	childID, err2 := strconv.Atoi(chi.URLParam(r, "childID"))
	amount, err3 := strconv.Atoi(r.FormValue("amount_paid"))
	if err1 != nil || err2 != nil || err3 != nil || amount < 0 {
		http.Error(w, "Nieprawidłowe dane", http.StatusBadRequest)
		return
	}

	if err := model.SetAmountPaid(h.DB, childID, eventID, amount); err != nil {
		http.Error(w, "Nie udało się zapisać płatności", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func (h *PaymentHandler) SetGoing(w http.ResponseWriter, r *http.Request) {
	eventID, err1 := strconv.Atoi(chi.URLParam(r, "eventID"))
	childID, err2 := strconv.Atoi(chi.URLParam(r, "childID"))
	if err1 != nil || err2 != nil {
		http.Error(w, "Nieprawidłowe dane", http.StatusBadRequest)
		return
	}

	user, _ := auth.GetUser(r)
	if owns, err := model.ChildBelongsToUser(h.DB, childID, user.ID); err != nil || !owns {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	going := r.FormValue("going") == "1"
	if err := model.SetGoing(h.DB, childID, eventID, going); err != nil {
		http.Error(w, "Nie udało się zapisać obecności", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func (h *PaymentHandler) TogglePaid(w http.ResponseWriter, r *http.Request) {
	eventID, err1 := strconv.Atoi(chi.URLParam(r, "eventID"))
	childID, err2 := strconv.Atoi(chi.URLParam(r, "childID"))
	if err1 != nil || err2 != nil {
		http.Error(w, "Nieprawidłowe dane", http.StatusBadRequest)
		return
	}

	user, _ := auth.GetUser(r)
	if owns, err := model.ChildBelongsToUser(h.DB, childID, user.ID); err != nil || !owns {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := model.TogglePaid(h.DB, childID, eventID); err != nil {
		http.Error(w, "Nie udało się zaktualizować płatności", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
