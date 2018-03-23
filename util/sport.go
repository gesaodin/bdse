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

	"github.com/tealeg/xlsx"
)

//LeerSport Sport
func (a *Archivo) LeerSport(ch chan []byte, tipo string) (bool, string) {
	fig := SParley
	posicionarchivo := 8
	if tipo == "f" {
		fig = SFigura
		posicionarchivo = 28
	}
	a.iniciarVariable(fig)

	insertar := a.Cabecera
	var coma string
	oid, b := a.CrearTraza(posicionarchivo, a.ConvertirTablaNumero(fig))
	if b != nil {
		m.Msj = "E# Sport17 : " + a.NombreDelArchivo + " " + b.Error()
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

		ree := strings.Replace(scan.Text(), ",", ".", -1)
		linea := strings.Fields(ree)
		l := len(linea)
		// fmt.Println(l)
		if l > 1 {
			if "VENDEDOR" == strings.Trim(linea[0], " ") {
				a.Leer = true
				a.CantidadLineas++
			}
			if a.Leer {
				if a.CantidadLineas > 2 && strings.Trim(linea[0], " ") != "TOTALES" && strings.Trim(linea[0], " ") != "" {
					coma = ","
				} else {
					coma = ""
				}
				insertar += coma
				if a.CantidadLineas > 1 && "TOTALES" != strings.Trim(linea[0], " ") && strings.Trim(linea[0], " ") != "" {

					re := regexp.MustCompile(`[(-)]`)
					agen := re.Split(linea[0], -1)
					p := strings.Replace(RComaXPunto(linea[5]), "-", "", -1)
					c := strings.Replace(RComaXPunto(linea[2]), "-", "", -1)
					agencia, venta := strings.Trim(agen[0], " "), RComaXPunto(linea[1])
					premio, comision := p, c
					insertar += "('" + agencia + "'," + venta + "," + premio + ","
					insertar += comision + ",1,'" + a.Fecha + "',Now(),8," + strconv.Itoa(oid) + ")"
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

//LeerSportXLSX Archivo en formato XLS 97-2003
func (a *Archivo) LeerSportXLSX(ch chan []byte, tipo string) (bool, string) {
	Exre := regexp.MustCompile(`[.()]`)
	ExArch := Exre.Split(a.NombreDelArchivo, -1)
	if ExArch[1] != "xlsx" {
		return a.LeerSport(ch, tipo)
	}
	fig := SParley
	posicionarchivo := 8
	if tipo == "f" {
		fig = SFigura
		posicionarchivo = 28
	}
	a.iniciarVariable(fig)

	insertar := a.Cabecera
	var coma string
	contar := 0
	oid, b := a.CrearTraza(posicionarchivo, a.ConvertirTablaNumero(fig))
	if b != nil {
		m.Msj = "E# Sport17 : " + a.NombreDelArchivo + " " + b.Error()
		m.Tipo = 33
		m.Tiempo = time.Now()
		j, _ := json.Marshal(m)
		ch <- j
		a.Canal <- j
		return false, ""
	}
	excelFileName := a.Ruta + a.NombreDelArchivo
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err)
	}
	for _, sheet := range xlFile.Sheets {

		for _, row := range sheet.Rows {

			var cel []string
			a.CantidadLineas++
			if a.CantidadLineas > 0 {
				contar++
				for _, cell := range row.Cells {
					text := cell.String()
					if strings.Trim(text, " ") != "" {
						cel = append(cel, text)
					}
				} //FIN DE LA CELDA

				l := len(cel)
				if l > 5 {
					if contar > 1 {
						coma = ","
					} else {
						coma = ""
					}
					re := regexp.MustCompile(`[-()]`)
					agen := re.Split(cel[1], -1)
					agencia, venta := strings.ToUpper(agen[0]), strings.Trim(cel[2], " ")
					premio, comision := strings.Trim(cel[6], " "), strings.Trim(cel[3], " ")
					insertar += coma
					insertar += "('" + agencia + "'," + venta + "," + premio + "," + comision
					insertar += ",1,'" + a.Fecha + "',Now()," + strconv.Itoa(posicionarchivo) + "," + strconv.Itoa(oid) + ")"
					a.Salvar = true
				}

			} //FIN DEL MAYOR A 7 FILAS

		} //FIN DE LA FILA
		// fmt.Println(insertar)
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
