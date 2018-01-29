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
func (a *Archivo) LeerPos(ch chan []byte, tipo string) (bool, string) {
	posicionarchivo, fig := TipoArchivo(tipo)
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

func TipoArchivo(tipo string) (posicionarchivo int, fig string) {

	switch tipo {
	case "1t":
		posicionarchivo = 2
		fig = "loteria"
		break
	case "2t":
		posicionarchivo = 3
		fig = "loteria"
		break
	case "3t":
		posicionarchivo = 4
		fig = "loteria"
		break
	case "1f":
		posicionarchivo = 15
		fig = "loteria"
		break
	case "2f":
		posicionarchivo = 16
		fig = "loteria"
		break
	case "3f":
		posicionarchivo = 17
		fig = "loteria"
		break
	default:
	}

	return
}
