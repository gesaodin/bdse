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
func (a *Archivo) LeerPos(ch chan []byte, tipo int) (bool, string) {
	a.iniciarVariable("loteria")
	insertar := a.Cabecera
	var coma string

	oid, b := a.CrearTraza(tipo, Loteria)
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
					insertar += strconv.Itoa(tipo) + "," + strconv.Itoa(oid) + ")"
					a.Salvar = true
				}
				a.CantidadLineas++
			}
		}
	}
	m.Tipo = 33
	m.Msj = "#" + a.NombreDelArchivo + " (" + strconv.Itoa(tipo) + ") Sin Registros"
	if a.Salvar {
		r, err := a.PostgreSQL.Exec(insertar)
		if err != nil {
			m.Tipo = 33
			m.Msj = "E#" + a.NombreDelArchivo + " (" + strconv.Itoa(tipo) + ") " + err.Error()
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
