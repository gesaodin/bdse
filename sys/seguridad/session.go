// seguridad (del latín securitas) cotidianamente se puede
// referir a la ausencia de riesgo o a la confianza en algo
// o en alguien. Sin embargo, el término puede tomar diversos
// sentidos según el área o campo a la que haga referencia en la
// seguridad. En términos generales, la seguridad se define como "el
// estado de bienestar que percibe y disfruta el ser humano".
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
		Domain:   "gokuserver",
		Path:     "/",
		MaxAge:   1800, //Media Hora en segundos
		HttpOnly: true,
	}
}

func (S *Session) Crear(w http.ResponseWriter, r *http.Request) {

}
