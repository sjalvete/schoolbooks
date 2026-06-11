package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"schoolbooks/internal/auth"
	"schoolbooks/internal/config"
	"schoolbooks/internal/db"
	"schoolbooks/internal/handler"
	"schoolbooks/internal/page"
	"schoolbooks/internal/templates"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/schoolbooks.db"
	}

	database := db.Init(dbPath)
	defer database.Close()

	cfg := config.Load()

	authHandler := &handler.AuthHandler{DB: database, Config: cfg}
	eventHandler := &handler.EventHandler{DB: database, Config: cfg}
	recipientHandler := &handler.RecipientHandler{DB: database, Config: cfg}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(auth.LoadUser)

	// public
	r.Get("/login", authHandler.ShowLogin)
	r.With(httprate.Limit(5, time.Minute, httprate.WithKeyFuncs(httprate.KeyByIP))).
		Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)

	// logged in
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			pd := page.New("Strona główna", r, w, cfg)
			templates.Home(pd).Render(r.Context(), w)
		})
		r.Get("/events", eventHandler.List)
	})

	// admin
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth)
		r.Use(auth.RequireAdmin)

		// r.Get("/users", userHandler.List)

		// r.Post("/users", userHandler.Create)
		// r.Get("/users/{id}", userHandler.EditUserForm)
		// r.Put("/users/{id}", userHandler.Update)

		r.Get("/events/manage", eventHandler.AdminList)
		r.Get("/events/new", eventHandler.NewEventForm)
		r.Get("/events/edit/{id}", eventHandler.EditEventForm)

		r.Post("/events", eventHandler.Create)
		r.Put("/events/{id}", eventHandler.Update)
		r.Delete("/events/{id}", eventHandler.Delete)

		r.Get("/recipients", recipientHandler.AdminList)
		r.Get("/recipients/new", recipientHandler.NewRecipientForm)
		r.Get("/recipients/edit/{id}", recipientHandler.EditRecipientForm)

		r.Post("/recipients", recipientHandler.Create)
		r.Put("/recipients/{id}", recipientHandler.Update)
		r.Delete("/recipients/{id}", recipientHandler.Delete)

	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
