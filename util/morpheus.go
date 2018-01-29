package util

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"
)

//LeerMorpheus Archivo de loteria, ch chan []byte
func (a *Archivo) LeerMorpheus(ch chan []byte) (bool, string) {
	a.iniciarVariable("loteria")
	insertar := a.Cabecera
	var coma string
	oid, b := a.CrearTraza(1, Loteria)
	if b != nil {
		m.Msj = "E# Morpheus: " + a.NombreDelArchivo + " " + b.Error()
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
		linea := strings.Split(ConvertirMonedaANumero(scan.Text()), "\t")
		if "CODIGO" == strings.Trim(linea[0], " ") {
			a.CantidadLineas++
			a.Leer = true
		}
		if a.Leer {
			if a.CantidadLineas > 2 && _Totales != strings.Trim(linea[0], " ") {
				coma = ","
			} else {
				coma = ""
			}
			insertar += coma
			if a.CantidadLineas > 1 && len(linea) == 8 && _Totales != strings.Trim(linea[0], " ") {
				agencia, venta := strings.Trim(linea[1], " "), strings.Trim(linea[3], " ")
				premio, comision := strings.Trim(linea[5], " "), strings.Trim(linea[7], " ")
				insertar += "('" + agencia + "'," + venta + "," + premio + ","
				insertar += comision + ",1,'" + a.Fecha + "',Now(),1," + strconv.Itoa(oid) + ")"
				a.Salvar = true
			}
			a.CantidadLineas++
		}
	}
	m.Tipo = 33
	m.Msj = "# Morpheus: " + a.NombreDelArchivo + " Sin Registros"
	if a.Salvar {
		r, err := a.PostgreSQL.Exec(insertar)
		if err != nil {
			m.Tipo = 33
			m.Msj = "E# Morpheus: " + a.NombreDelArchivo + " " + err.Error()
		} else {
			m.Tipo = 1
			i, _ := r.RowsAffected()
			a.Registros = int(i)
			m.Msj = "Morpheus se registrarón: " + strconv.Itoa(a.Registros) + " Filas."

		}
	}
	a.ModificarTraza()
	m.Tiempo = time.Now()
	j, _ := json.Marshal(m)
	ch <- j
	a.Canal <- j
	return true, insertar
}
