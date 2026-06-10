package page

import (
	"net/http"
	"schoolbooks/internal/auth"
	"schoolbooks/internal/session"
)

type Data struct {
	Title   string
	User    *auth.SessionUser
	Flash   *session.Flash
	IsAdmin bool
}

func New(title string, r *http.Request, w http.ResponseWriter) *Data {
	user, ok := auth.GetUser(r)
	d := &Data{
		Title: title,
		Flash: session.PopFlash(w, r),
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
