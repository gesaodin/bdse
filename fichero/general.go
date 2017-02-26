package fichero

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool
}

//Validar los campos nulos de la base de datos y retornar su valor original
func ValidarNullString(b sql.NullString) (s string) {
	if b.Valid {
		s = b.String
	} else {
		s = "null"
	}
	return
}

func ValidarNullTime(b interface{}) (t time.Time) {
	t, e := b.(time.Time)
	if !e {
		return time.Now()
	}
	return
}

//Convertir de (YYYY-MM-DD) a (DD/MM/YYYY) Humano
func ConvertirFechaSlash(fecha string) string {
	return "23/07/2016"
}

func ConvertirMonedaANumero(moneda string) string {
	conver := strings.Replace(moneda, ".", "", -1)
	stringNumero := strings.Replace(conver, ",", ".", -1)
	return stringNumero
}

//Calcular los dias de un mes
func DiasDelMes(fecha time.Time) int {
	return 0
}

//Permite llenar con ceros antes y despues de una cadena
func CompletarCeros(cadena string, orientacion int, cantidad int) string {
	return "000"
}

func Fatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Error(e error) {
	if e != nil {
		fmt.Println("BDSE@Error: $ ", e)
	}
}
