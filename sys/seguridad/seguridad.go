// Seguridad (del latín securitas)1 cotidianamente se puede
// referir a la ausencia de riesgo o a la confianza en algo
// o en alguien. Sin embargo, el término puede tomar diversos
// sentidos según el área o campo a la que haga referencia en la
// seguridad. En términos generales, la seguridad se define como "el
// estado de bienestar que percibe y disfruta el ser humano".
package seguridad

import (
	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	Encriptamiento             = "md5"
	ActivarLimiteDeConexion    = true
	DesactivarLimiteDeConexion = false
)

var (
	LlavePrivada *rsa.PrivateKey
	LlavePublica *rsa.PublicKey
	LlaveJWT     string
)

func init() {
	mySigningKey := []byte("AllYourBase")
	claims := Reclamaciones{
		"admin",
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	//fmt.Printf("%v %v", ss, err)
	LlaveJWT = ss

}
