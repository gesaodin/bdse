package util

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//LeerCyberParley Cyber Parley
func (a *Archivo) LeerCyberParley(ch chan []byte) (bool, string) {
	a.iniciarVariable("parley")
	insertar := a.Cabecera
	var coma string

	oid, b := a.CrearTraza(7, Parley)
	if b != nil {
		m.Msj = "E# CyberParley: " + a.NombreDelArchivo + " " + b.Error()
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
		linea := strings.Split(ConvertirMonedaANumero(scan.Text()), ";")
		//l := len(linea)

		if strings.Trim(linea[0], " ") == "Tipo Entidad" {
			a.Leer = true
			a.CantidadLineas++
		}
		if a.Leer {
			if a.CantidadLineas > 2 && strings.Trim(linea[0], " ") == "Agencia" {
				coma = ","
			} else {
				coma = ""
			}
			insertar += coma
			if a.CantidadLineas > 1 && strings.Trim(linea[0], " ") == "Agencia" {

				re := regexp.MustCompile(`[(-)]`)
				agen := re.Split(linea[1], -1)
				// fmt.Println(agen[2])

				c := strings.Replace(strings.Trim(linea[4], " "), "-", "", -1)
				p := strings.Replace(strings.Trim(linea[5], " "), "-", "", -1)
				agencia, venta := strings.Trim(agen[0], " "), strings.Trim(linea[3], " ")
				premio, comision := p, c
				insertar += "('" + agencia + "'," + venta + "," + premio + ","
				insertar += comision + ",1,'" + a.Fecha + "',Now(),7," + strconv.Itoa(oid) + ")"
				a.Salvar = true
			}
			a.CantidadLineas++
		}

	}

	m.Tipo = 33
	m.Msj = "E#" + a.NombreDelArchivo + " Sin Registros"
	if a.Salvar {
		r, err := a.PostgreSQL.Exec(insertar)
		if err != nil {
			m.Tipo = 33
			m.Msj = "E#" + a.NombreDelArchivo + err.Error()
		} else {
			m.Tipo = 1
			i, _ := r.RowsAffected()
			a.Registros = int(i)
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
