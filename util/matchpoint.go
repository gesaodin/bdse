package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

//LeerMatchPoint Archivo en formato XLS 97-2003
func (a *Archivo) LeerMatchPoint(ch chan []byte, tipo string) (bool, string) {

	fig := SParley
	posicionarchivo := 27
	if tipo == "f" {
		fig = SFigura
		posicionarchivo = 29
	}
	a.iniciarVariable(fig)

	insertar := a.Cabecera
	var coma string
	contar := 0
	oid, b := a.CrearTraza(posicionarchivo, Loteria)
	if b != nil {
		m.Msj = "E# Maticlot : " + a.NombreDelArchivo + " " + b.Error()
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
			if a.CantidadLineas > 8 {
				contar++
				for _, cell := range row.Cells {
					text := cell.String()
					if strings.Trim(text, " ") != "" {
						cel = append(cel, text)
					}
				} //FIN DE LA CELDA

				l := len(cel)
				if l > 6 {
					if contar > 1 && strings.Trim(cel[0], " ") != "Total" {
						coma = ","
					} else {
						coma = ""
					}
					agencia, venta := strings.ToUpper(cel[2]), strings.Trim(cel[3], " ")
					premio, comision := strings.Trim(cel[4], " "), strings.Trim(cel[5], " ")
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
