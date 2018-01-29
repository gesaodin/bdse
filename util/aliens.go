package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//LeerPos Archivos de Loteria
func (a *Archivo) LeerAliens(ch chan []byte, tipo string) (bool, string) {
	posicionarchivo, fig := TipoArchivoAliens(tipo)
	a.iniciarVariable(fig)
	insertar := a.Cabecera
	var coma string

	oid, b := a.CrearTraza(posicionarchivo, Loteria)
	if b != nil {
		m.Msj = "E# : " + a.NombreDelArchivo + " " + b.Error()
		m.Tipo = 33
		m.Tiempo = time.Now()
		j, _ := json.Marshal(m)
		ch <- j
		a.Canal <- j
		return false, ""
	}
	archivo, err := os.Open(a.Ruta + a.NombreDelArchivo)
	Error(err)
	scan := bufio.NewScanner(archivo)
	for scan.Scan() {

		linea := strings.Fields(scan.Text())
		l := len(linea)
		//fmt.Println("# ", l, linea)
		if l > 3 {
			if "TAQUILLA" == strings.Trim(linea[0], " ") {
				a.Leer = true
				a.CantidadLineas++
			}
			if a.Leer {
				if a.CantidadLineas > 2 && strings.Trim(linea[0], " ") != "TOTALES:" && strings.Trim(linea[0], " ") != "" {
					coma = ","
				} else {
					coma = ""
				}
				insertar += coma
				if a.CantidadLineas > 1 && "TOTALES:" != strings.Trim(linea[0], " ") && strings.Trim(linea[0], " ") != "" {

					re := regexp.MustCompile(`[-()]`)
					agen := re.Split(linea[0], -1)
					agencia, venta := agen[1], strings.Trim(linea[l-3], " ")
					premio, comision := strings.Trim(linea[l-1], " "), strings.Trim(linea[l-2], " ")
					insertar += "('" + agencia + "'," + venta + "," + premio + ","
					insertar += comision + ",1,'" + a.Fecha + "',Now(),"
					insertar += strconv.Itoa(posicionarchivo) + "," + strconv.Itoa(oid) + ")"
					a.Salvar = true
				}
				a.CantidadLineas++
			}
		}
	}
	m.Tipo = 33
	m.Msj = "#" + a.NombreDelArchivo + " (" + strconv.Itoa(posicionarchivo) + ") Sin Registros"
	if a.Salvar {
		r, err := a.PostgreSQL.Exec(insertar)
		if err != nil {
			m.Tipo = 33
			m.Msj = "E#" + a.NombreDelArchivo + " (" + strconv.Itoa(posicionarchivo) + ") " + err.Error()
			fmt.Println(m.Msj)
		} else {
			m.Tipo = 1
			i, _ := r.RowsAffected()
			a.Registros = int(i)
			// m.Msj = a.NombreDelArchivo + " (" + strconv.Itoa(tipo) + ") se registrarón: " + filas + " Filas."
			m.Msj = a.NombreDelArchivo + " se registrarón: " + strconv.Itoa(a.Registros) + " Filas."

		}
	}
	a.ModificarTraza()
	m.Tiempo.Format("2006-01-01 00:00:00")
	m.Tiempo = time.Now()
	j, _ := json.Marshal(m)
	ch <- j
	a.Canal <- j
	return true, insertar
}

func TipoArchivoAliens(tipo string) (posicionarchivo int, fig string) {

	switch tipo {
	case "1t":
		posicionarchivo = 9
		fig = "loteria"
		break
	case "2t":
		posicionarchivo = 10
		fig = "loteria"
		break
	case "3t":
		posicionarchivo = 11
		fig = "loteria"
		break
	case "1f":
		posicionarchivo = 15
		fig = "figura"
		break
	case "2f":
		posicionarchivo = 16
		fig = "figura"
		break
	case "3f":
		posicionarchivo = 17
		fig = "figura"
		break
	case "1c":
		posicionarchivo = 15
		fig = "truco"
		break
	case "2c":
		posicionarchivo = 16
		fig = "truco"
		break
	case "3c":
		posicionarchivo = 17
		fig = "truco"
		break
	case "1p":
		posicionarchivo = 15
		fig = "pescalo"
		break
	case "2p":
		posicionarchivo = 16
		fig = "pescalo"
		break
	case "3p":
		posicionarchivo = 17
		fig = "pescalo"
		break

	default:
	}

	return
}
