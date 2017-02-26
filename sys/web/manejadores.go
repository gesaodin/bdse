package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gesaodin/bdse/sys/seguridad"
	"github.com/gesaodin/bdse/util"
)

func Cabecera(w http.ResponseWriter, origen string) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	w.Header().Set("Access-Control-Allow-Origin", origen)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

}

func CabeceraRechazada(w http.ResponseWriter, estatus int, m string) {
	w.WriteHeader(estatus)
	msj := []byte(m)
	w.Write(msj)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var u seguridad.Usuario

	Cabecera(w, r.Header.Get("Origin"))
	e := json.NewDecoder(r.Body).Decode(&u)
	util.Error(e)

	fmt.Println("Pasando el Decode", u)
	if u.Nombre == "carlos" && u.Clave == "za63qj2p" {
		u.Nombre = "Carlos"
		u.Clave = ""
		u.Id = 0
		token := seguridad.GenerarJWT(u)
		result := seguridad.RespuestaToken{Token: token}
		j, e := json.Marshal(result)
		util.Error(e)
		Mensajeria.Usuario["gpanel"].ch <- []byte("Iniciando Sesión")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	} else {
		w.Header().Set("Content-Type", "application/text")
		fmt.Println("Error en la conexion del usuario")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "Usuario y clave no validas")
	}
}

func ValidarToken(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := seguridad.Stores.Get(r, "session-bdse")
		fmt.Println("Conexion establecida desde: ", r.Header.Get("Origin"))
		Cabecera(w, r.Header.Get("Origin"))
		token, e := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &seguridad.Reclamaciones{}, func(token *jwt.Token) (interface{}, error) {
			return seguridad.LlavePublica, nil
		})

		if e != nil {
			switch e.(type) {
			case *jwt.ValidationError:
				vErr := e.(*jwt.ValidationError)
				switch vErr.Errors {
				case jwt.ValidationErrorExpired:
					fmt.Println("Expirar")
					w.WriteHeader(http.StatusUnauthorized)
					msj := []byte("El token ha expirado")
					w.Write(msj)
					return
				case jwt.ValidationErrorSignatureInvalid:
					w.WriteHeader(http.StatusForbidden)
					msj := []byte("La firma del token no coincide")
					w.Write(msj)
					return
				default:
					w.WriteHeader(http.StatusForbidden)
					msj := []byte("Token invalido")
					w.Write(msj)
					return
				}
			default:
				fmt.Fprintln(w, "El token no es valido")
				return
			}
		}
		fmt.Println("Validando")
		if token.Valid {
			session.Values["ok"] = true
			session.Values["name"] = ""
			session.Save(r, w)
			fn(w, r)
		} else {
			CabeceraRechazada(w, http.StatusForbidden, "El token no es valido")
			return
		}
	})
}
