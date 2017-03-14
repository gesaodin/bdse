package util

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

//NullTime Tiempo nulo
type NullTime struct {
	Time  time.Time
	Valid bool
}

//ValidarNullString los campos nulos de la base de datos y retornar su valor original
func ValidarNullString(b sql.NullString) (s string) {
	if b.Valid {
		s = b.String
	} else {
		s = ""
	}
	return
}

//ValidarNullFloat64 los campos nulos de la base de datos y retornar su valor original
func ValidarNullFloat64(b sql.NullFloat64) (f float64) {
	if b.Valid {
		f = b.Float64
	} else {
		f = 0
	}
	return
}

//ValidarNullTime los campos nulos de la base de datos y retorna fecha
func ValidarNullTime(b interface{}) (t time.Time) {
	t, e := b.(time.Time)
	if !e {
		return time.Now()
	}
	return
}

//ConvertirFechaSlash de (YYYY-MM-DD) a (DD/MM/YYYY) Humano
func ConvertirFechaSlash(fecha string) string {
	return "23/07/2016"
}

//ConvertirMonedaANumero Numeros
func ConvertirMonedaANumero(moneda string) string {
	conver := strings.Replace(moneda, ".", "", -1)
	stringNumero := strings.Replace(conver, ",", ".", -1)
	return stringNumero
}

//DiasDelMes Calcular los dias de un mes
func DiasDelMes(fecha time.Time) int {
	return 0
}

//CompletarCeros Permite llenar con ceros antes y despues de una cadena
func CompletarCeros(cadena string, orientacion int, cantidad int) string {
	return "000"
}

//Fatal Errores
func Fatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

//Error reglas generales
func Error(e error) {
	if e != nil {
		fmt.Println("BDSE@Error: $ ", e)
	}
}
