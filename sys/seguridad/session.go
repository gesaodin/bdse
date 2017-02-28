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
		Domain:   "200.82.247.127",
		Path:     "/",
		MaxAge:   1800, //Media Hora en segundos
		HttpOnly: true,
	}
}

func (S *Session) Crear(w http.ResponseWriter, r *http.Request) {

}
