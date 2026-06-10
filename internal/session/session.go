package session

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("change-this-secret-before-production"))

type Flash struct {
	Message string
	Type    string
}

func init() {
	gob.Register(int(0))
	gob.Register(int64(0))
	gob.Register("")
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return Store.Get(r, "session")
}

func GetUserID(r *http.Request) int {
	session, err := Store.Get(r, "session")
	if err != nil {
		return 0
	}
	id, _ := session.Values["user_id"].(int)
	return id
}

func GetUserName(r *http.Request) string {
	session, err := Store.Get(r, "session")
	if err != nil {
		return ""
	}
	name, _ := session.Values["name"].(string)
	return name
}

func GetUserRole(r *http.Request) string {
	session, err := Store.Get(r, "session")
	if err != nil {
		return ""
	}
	role, _ := session.Values["role"].(string)
	return role
}

func IsAdmin(r *http.Request) bool {
	return GetUserRole(r) == "su"
}

func SetFlash(w http.ResponseWriter, r *http.Request, message, flashType string) {
	session, _ := Store.Get(r, "session")
	session.Values["flash_message"] = message
	session.Values["flash_type"] = flashType
	session.Save(r, w)
}

func PopFlash(w http.ResponseWriter, r *http.Request) *Flash {
	session, _ := Store.Get(r, "session")
	msg, ok := session.Values["flash_message"].(string)
	if !ok || msg == "" {
		return nil
	}
	flashType, _ := session.Values["flash_type"].(string)

	// delete after reading
	delete(session.Values, "flash_message")
	delete(session.Values, "flash_type")
	session.Save(r, w)

	return &Flash{Message: msg, Type: flashType}
}
