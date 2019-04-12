package util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

//LeerMatrix Archivo en formato XLS 97-2003
func (a *Archivo) LeerMatrix(ch chan []byte, tipo string) (bool, string) {

	fig := SFigura
	posicionarchivo := 26
	a.iniciarVariable(fig)

	insertar := a.Cabecera
	var coma string
	contar := 0
	oid, b := a.CrearTraza(posicionarchivo, a.ConvertirTablaNumero(fig))
	if b != nil {
		m.Msj = "E# Matrix : " + a.NombreDelArchivo + " " + b.Error()
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
		fmt.Println("Error en el formato del excel ( *.xlsx ): ", err.Error())
	}
	for _, sheet := range xlFile.Sheets {

		for _, row := range sheet.Rows {

			var cel []string
			a.CantidadLineas++
			if a.CantidadLineas > 3 {
				contar++
				for _, cell := range row.Cells {
					text := cell.String()
					if strings.Trim(text, " ") != "" {
						cel = append(cel, text)
					}
				} //FIN DE LA CELDA
				l := len(cel)

				if l == 8 {
					if contar > 2 {
						coma = ","
					} else {
						coma = ""
					}
					re := regexp.MustCompile(`[-()]`)
					agen := re.Split(cel[1], -1)

					agencia, venta := strings.ToUpper(agen[0]), RComaXPunto(cel[2])
					premio, comision := RComaXPunto(cel[4]), RComaXPunto(cel[3])
					insertar += coma
					insertar += "('" + agencia + "'," + venta + "," + premio + "," + comision
					insertar += ",1,'" + a.Fecha + "',Now()," + strconv.Itoa(posicionarchivo) + "," + strconv.Itoa(oid) + ")"
					a.Salvar = true
				}

			} //FIN DEL MAYOR A 7 FILAS

		} //FIN DE LA FILA
	}
	//fmt.Println(insertar)
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
