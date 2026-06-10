package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"schoolbooks/internal/model"
	"schoolbooks/internal/page"
	"schoolbooks/internal/templates"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	DB *sql.DB
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	pd := page.New("Events", r, w)
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
		http.Error(w, "could not load events", http.StatusInternalServerError)
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
		fmt.Sprintf("%s %d", time.Month(month).String(), year),
	).Render(r.Context(), w)
}

func (h *EventHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	pd := page.New("Events", r, w)
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
		http.Error(w, "could not load events", http.StatusInternalServerError)
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
		fmt.Sprintf("%s %d", time.Month(month).String(), year),
	).Render(r.Context(), w)
}

func (h *EventHandler) NewEventForm(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")

	templates.NewEventForm(date).Render(r.Context(), w)
}

func (h *EventHandler) EditEventForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	e, err := model.GetEventByID(h.DB, id)

	if err != nil {
		http.Error(w, "could not load event", http.StatusInternalServerError)
		return
	}

	templates.EditEventForm(e).Render(r.Context(), w)
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	date := r.FormValue("date")
	price := r.FormValue("price")

	if err := model.CreateEvent(h.DB, title, description, price, date); err != nil {
		http.Error(w, "could not create event", http.StatusInternalServerError)
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

	if err := model.UpdateEvent(h.DB, id, title, description, price, date); err != nil {
		http.Error(w, "could not update event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := model.DeleteEvent(h.DB, id); err != nil {
		http.Error(w, "could not delete event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
