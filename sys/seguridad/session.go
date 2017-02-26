package seguridad

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Session struct {
	Nombre string
	Acceso string
	Nivel  int
}

var Stores = sessions.NewCookieStore([]byte("#za63qj2p-6pt33pSUz#"))

func init() {
	Stores.Options = &sessions.Options{
		Domain:   "192.168.1.100",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}
}

func (S *Session) Crear(w http.ResponseWriter, r *http.Request) {

}
