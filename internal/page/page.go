package page

import (
	"log"
	"net/http"
	"schoolbooks/internal/auth"
	"schoolbooks/internal/config"
	"schoolbooks/internal/session"
)

type Data struct {
	Title   string
	User    *auth.SessionUser
	Flash   *session.Flash
	Config  *config.Config
	IsAdmin bool
}

func New(title string, r *http.Request, w http.ResponseWriter, c *config.Config) *Data {
	log.Printf("debug mode: %v", c.Debug)
	user, ok := auth.GetUser(r)
	d := &Data{
		Title:  title,
		Flash:  session.PopFlash(w, r),
		Config: c,
	}
	if ok {
		d.User = &user
		d.IsAdmin = user.Role == "admin"

	}
	return d
}

func (d *Data) SetFlash(message, flashType string) {
	d.Flash = &session.Flash{Message: message, Type: flashType}
}
