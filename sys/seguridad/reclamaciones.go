package seguridad

import jwt "github.com/dgrijalva/jwt-go"

type Reclamaciones struct {
	Rol string
	jwt.StandardClaims
}
