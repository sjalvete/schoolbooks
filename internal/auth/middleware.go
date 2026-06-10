package auth

import (
	"context"
	"log"
	"net/http"
	"schoolbooks/internal/session"
	"time"
)

type contextKey string

const UserKey contextKey = "user"

type SessionUser struct {
	ID   int
	Name string
	Role string
}

func LoadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, _ := session.GetSession(r)
		if last, ok := s.Values["last_active"].(int64); ok {
			if time.Now().Unix()-last > 15*60 {
				s.Options.MaxAge = -1
				s.Save(r, w)
				http.Redirect(w, r, "/login?expired=1", http.StatusSeeOther)
				return
			}
		}

		if id, ok := s.Values["user_id"].(int); ok && id != 0 {
			s.Values["last_active"] = time.Now().Unix()

			if err := s.Save(r, w); err != nil {
				log.Printf("session save error: %v", err)
			}

			user := SessionUser{
				ID:   id,
				Name: s.Values["name"].(string),
				Role: s.Values["role"].(string),
			}
			ctx := context.WithValue(r.Context(), UserKey, user)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func GetUser(r *http.Request) (SessionUser, bool) {
	user, ok := r.Context().Value(UserKey).(SessionUser)
	return user, ok
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetUser(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUser(r)
		if !ok || user.Role != "admin" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
