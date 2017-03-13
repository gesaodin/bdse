// se refiere a la confianza en algo.
package seguridad

import (
	"crypto/rsa"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gesaodin/bdse/util"
)

//Constantes Generales
const (
	Encriptamiento             = "md5"
	ActivarLimiteDeConexion    = true
	DesactivarLimiteDeConexion = false
)

//Variables de Seguridad
var (
	LlavePrivada *rsa.PrivateKey
	LlavePublica *rsa.PublicKey
	LlaveJWT     string
)

//init Función inicial del sistema
func init() {
	bytePrivados, err := ioutil.ReadFile("./sys/seguridad/private.rsa")
	util.Fatal(err)
	LlavePrivada, err = jwt.ParseRSAPrivateKeyFromPEM(bytePrivados)
	bytePublicos, err := ioutil.ReadFile("./sys/seguridad/public.rsa.pub")
	util.Fatal(err)
	LlavePublica, err = jwt.ParseRSAPublicKeyFromPEM(bytePublicos)
}

//GenerarJWT Json Web Token
func GenerarJWT(u Usuario) string {
	peticion := Reclamaciones{
		Usuario: u,
		Rol:     "Development",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "Conexion Bus Empresarial",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, peticion)
	rs, e := token.SignedString(LlavePrivada)
	util.Fatal(e)
	return rs
}

// func Cabecera(w http.ResponseWriter, origen string) {
// 	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
//
// 	w.Header().Set("Access-Control-Allow-Origin", origen)
// 	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
// 	w.Header().Set("Access-Control-Allow-Credentials", "true")
//
// }
//
// func CabeceraRechazada(w http.ResponseWriter, estatus int, m string) {
// 	w.WriteHeader(estatus)
// 	msj := []byte(m)
// 	w.Write(msj)
// }
//
// func Login(w http.ResponseWriter, r *http.Request) {
// 	var u Usuario
//
// 	Cabecera(w, r.Header.Get("Origin"))
// 	e := json.NewDecoder(r.Body).Decode(&u)
// 	util.Error(e)
// 	fmt.Println("Pasando el Decode")
// 	if u.Nombre == "carlos" && u.Clave == "za63qj2p" {
// 		u.Nombre = "Carlos"
// 		u.Clave = ""
// 		u.Id = 0
// 		token := GenerarJWT(u)
// 		result := RespuestaToken{token}
// 		j, e := json.Marshal(result)
// 		util.Error(e)
//
// 		w.WriteHeader(http.StatusOK)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(j)
// 	} else {
// 		w.Header().Set("Content-Type", "application/text")
// 		fmt.Println("Error en la conexion del usuario")
// 		w.WriteHeader(http.StatusForbidden)
// 		fmt.Fprintln(w, "Usuario y clave no validas")
// 	}
// }
// func LoginToken(w http.ResponseWriter, r *http.Request) {
// 	var u Usuario
//
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	session, _ := Stores.Get(r, "session-bdse")
// 	e := json.NewDecoder(r.Body).Decode(&u)
// 	util.Error(e)
// 	fmt.Println("Pasando el Decode")
// 	if u.Nombre == "carlos" && u.Clave == "za63qj2p" {
// 		u.Nombre = "Carlos"
// 		u.Clave = ""
// 		u.Id = 0
// 		token := GenerarJWT(u)
// 		result := RespuestaToken{token}
// 		j, e := json.Marshal(result)
// 		util.Error(e)
//
// 		w.WriteHeader(http.StatusOK)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(j)
// 		session.Values["ok"] = true
// 		session.Values["name"] = u.Nombre
// 		session.Save(r, w)
// 		fmt.Println("Creando Session")
// 	} else {
// 		w.Header().Set("Content-Type", "application/text")
// 		session.Values["ok"] = false
// 		session.Values["name"] = "Desconocido"
// 		session.Save(r, w)
// 		fmt.Println("Error en la conexion del usuario")
// 		w.WriteHeader(http.StatusForbidden)
// 		fmt.Fprintln(w, "Usuario y clave no validas")
// 	}
// }

// func ValidarTokenNew(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
// 	origin := r.Header.Get("Origin")
// 	w.Header().Set("Access-Control-Allow-Origin", origin)
// 	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
// 	w.Header().Set("Access-Control-Allow-Credentials", "true")
//
// 	fmt.Println("Entrando")
// 	token, e := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
// 		return LlavePublica, nil
// 	})
// 	if e != nil {
// 		switch e.(type) {
// 		case *jwt.ValidationError:
// 			vErr := e.(*jwt.ValidationError)
// 			switch vErr.Errors {
// 			case jwt.ValidationErrorExpired:
// 				w.WriteHeader(http.StatusForbidden)
// 				msj := []byte("El token ha expirado")
// 				w.Write(msj)
// 				return
// 			case jwt.ValidationErrorSignatureInvalid:
// 				fmt.Fprintln(w, "('La firma del token no coincide')")
// 				return
// 			default:
// 				fmt.Fprintln(w, "('Token invalido')")
// 				return
// 			}
// 		default:
// 			fmt.Fprintln(w, "('El token no es valido')")
// 			return
// 		}
// 	}
// 	fmt.Println("Validando")
//
// 	if token.Valid {
// 		w.WriteHeader(http.StatusAccepted)
// 		msj := "Bienvenido"
// 		fmt.Fprintln(w, msj)
// 	} else {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		msj := "Su token no es valido"
// 		fmt.Fprintln(w, msj)
// 	}
// }

// func ValidarToken(fn http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		session, _ := Stores.Get(r, "session-bdse")
// 		fmt.Println("Conexion establecida desde: ", r.Header.Get("Origin"))
// 		Cabecera(w, r.Header.Get("Origin"))
// 		token, e := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &Reclamaciones{}, func(token *jwt.Token) (interface{}, error) {
// 			return LlavePublica, nil
// 		})
//
// 		if e != nil {
// 			switch e.(type) {
// 			case *jwt.ValidationError:
// 				vErr := e.(*jwt.ValidationError)
// 				switch vErr.Errors {
// 				case jwt.ValidationErrorExpired:
// 					fmt.Println("Expirar")
// 					w.WriteHeader(http.StatusUnauthorized)
// 					msj := []byte("El token ha expirado")
// 					w.Write(msj)
// 					return
// 				case jwt.ValidationErrorSignatureInvalid:
// 					w.WriteHeader(http.StatusForbidden)
// 					msj := []byte("La firma del token no coincide")
// 					w.Write(msj)
// 					return
// 				default:
// 					w.WriteHeader(http.StatusForbidden)
// 					msj := []byte("Token invalido")
// 					w.Write(msj)
// 					return
// 				}
// 			default:
// 				fmt.Fprintln(w, "El token no es valido")
// 				return
// 			}
// 		}
// 		fmt.Println("Validando")
// 		if token.Valid {
// 			session.Values["ok"] = true
// 			session.Values["name"] = ""
// 			session.Save(r, w)
// 			fn(w, r)
// 		} else {
// 			CabeceraRechazada(w, http.StatusForbidden, "El token no es valido")
// 			return
// 		}
// 	})
// }
