package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"schoolbooks/internal/config"
	"schoolbooks/internal/locale"
	"schoolbooks/internal/model"
	"schoolbooks/internal/page"
	"schoolbooks/internal/templates"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	DB     *sql.DB
	Config *config.Config
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	pd := page.New("Wydarzenia", r, w, h.Config)
	now := time.Now()
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	view := r.URL.Query().Get("view")

	var events []model.Event
	var err error

	if year == 0 {
		year = now.Year()
	}
	if month == 0 {
		month = int(now.Month())
	}
	if view == "" {
		view = "calendar"
	}

	if view == "agenda" {
		events, err = model.ListFutureEvents(h.DB)
	} else {
		events, err = model.ListEventsByMonth(h.DB, year, month)
	}

	if err != nil {
		http.Error(w, "Nie udało się załadować wydarzeń", http.StatusInternalServerError)
		return
	}

	eventMap := make(map[int][]model.Event)
	for _, e := range events {
		t, _ := time.Parse("2006-01-02", e.Date[:10])
		eventMap[t.Day()] = append(eventMap[t.Day()], e)
	}

	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.Local).Day()

	templates.UserEvents(
		pd,
		year, month, view,
		firstDay.Weekday(),
		daysInMonth,
		eventMap,
		events,
		fmt.Sprintf("%s %d", locale.PolishMonth(time.Month(month)), year),
	).Render(r.Context(), w)
}

func (h *EventHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	pd := page.New("Wydarzenia", r, w, h.Config)
	now := time.Now()
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	view := r.URL.Query().Get("view")

	var events []model.Event
	var err error

	if year == 0 {
		year = now.Year()
	}
	if month == 0 {
		month = int(now.Month())
	}
	if view == "" {
		view = "calendar"
	}

	if view == "agenda" {
		events, err = model.ListFutureEvents(h.DB)
	} else {
		events, err = model.ListEventsByMonth(h.DB, year, month)
	}

	if err != nil {
		http.Error(w, "Nie udało się załadować wydarzeń", http.StatusInternalServerError)
		return
	}

	eventMap := make(map[int][]model.Event)
	for _, e := range events {
		t, _ := time.Parse("2006-01-02", e.Date[:10])
		eventMap[t.Day()] = append(eventMap[t.Day()], e)
	}

	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.Local).Day()

	templates.AdminEvents(
		pd,
		year, month, view,
		firstDay.Weekday(),
		daysInMonth,
		eventMap,
		events,
		fmt.Sprintf("%s %d", locale.PolishMonth(time.Month(month)), year),
	).Render(r.Context(), w)
}

func (h *EventHandler) NewEventForm(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")

	recipients, err := model.ListRecipients(h.DB)
	if err != nil {
		http.Error(w, "Nie udało się załadować odbiorców", http.StatusInternalServerError)
		return
	}

	templates.NewEventForm(date, recipients).Render(r.Context(), w)
}

func (h *EventHandler) EditEventForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	e, err := model.GetEventByID(h.DB, id)

	if err != nil {
		http.Error(w, "Nie udało się załadować wydarzenia", http.StatusInternalServerError)
		return
	}

	recipients, err := model.ListRecipients(h.DB)
	if err != nil {
		http.Error(w, "Nie udało się załadować odbiorców", http.StatusInternalServerError)
		return
	}

	templates.EditEventForm(e, recipients).Render(r.Context(), w)
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	date := r.FormValue("date")
	price := r.FormValue("price")
	skippable := r.FormValue("skippable") != ""
	recipientID := parseRecipientID(r.FormValue("recipient_id"))

	if err := model.CreateEvent(h.DB, title, description, price, date, skippable, recipientID); err != nil {
		http.Error(w, "Nie udało się utworzyć wydarzenia", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	title := r.FormValue("title")
	description := r.FormValue("description")
	date := r.FormValue("date")
	price := r.FormValue("price")
	skippable := r.FormValue("skippable") != ""
	recipientID := parseRecipientID(r.FormValue("recipient_id"))

	if err := model.UpdateEvent(h.DB, id, title, description, price, date, skippable, recipientID); err != nil {
		http.Error(w, "Nie udało się zaktualizować wydarzenia", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func parseRecipientID(v string) *int {
	id, err := strconv.Atoi(v)
	if err != nil {
		return nil
	}
	return &id
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := model.DeleteEvent(h.DB, id); err != nil {
		http.Error(w, "Nie udało się usunąć wydarzenia", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
