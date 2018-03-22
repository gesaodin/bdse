package util

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"
)

//LeerIlbanquero Consultar en ilbanquero con el proveedor
func (a *Archivo) LeerIlbanquero(ch chan []byte, tipo string) (bool, string) {
	posicionarchivo, fig := tipoIlbanquero(tipo)
	a.iniciarVariable(fig)

	insertar := a.Cabecera
	var coma string
	oid, b := a.CrearTraza(posicionarchivo, a.ConvertirTablaNumero(fig))
	if b != nil {
		m.Msj = "E# Ilbanquero: " + a.NombreDelArchivo + " " + b.Error()
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

		l := len(linea)
		if l > 11 && strings.Trim(linea[1], " ") != "0.00" {
			if "Agentes" == strings.Trim(linea[0], " ") {
				a.Leer = true
				a.CantidadLineas++
			}
			if a.Leer {
				if a.CantidadLineas > 2 && strings.Trim(linea[0], " ") != "Total" && strings.Trim(linea[0], " ") != "" {
					coma = ","
				} else {
					coma = ""
				}
				insertar += coma
				if a.CantidadLineas > 1 && "Total" != strings.Trim(linea[0], " ") && strings.Trim(linea[0], " ") != "" {
					agencia, venta := strings.Trim(linea[0], " "), strings.Trim(linea[1], " ")
					premio, comision := strings.Trim(linea[5], " "), strings.Trim(linea[7], " ")
					insertar += "('" + agencia + "'," + venta + "," + premio + "," + comision
					insertar += ",1,'" + a.Fecha + "',Now()," + strconv.Itoa(posicionarchivo) + "," + strconv.Itoa(oid) + ")"
					a.Salvar = true
				}
				a.CantidadLineas++
			}
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

func tipoIlbanquero(tipo string) (posicionarchivo int, fig string) {
	switch tipo {
	case "t":
		posicionarchivo = 23
		fig = SLoteria
		break
	case "p":
		posicionarchivo = 6
		fig = SParley
		break
	case "f":
		posicionarchivo = 24
		fig = SFigura
		break
	default:
		break
	}
	return
}
