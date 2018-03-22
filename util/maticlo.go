package util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
)

//LeerMaticlo Archivo en formato XLS 97-2003
func (a *Archivo) LeerMaticlo(ch chan []byte, tipo string) (bool, string) {

	fig := SLoteria
	posicionarchivo := 5
	if tipo == "f" {
		fig = SFigura
		posicionarchivo = 13
	}
	a.iniciarVariable(fig)

	insertar := a.Cabecera
	var coma string
	contar := 0
	oid, b := a.CrearTraza(posicionarchivo, a.ConvertirTablaNumero(fig))
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
	if xlFile, err := xls.Open(excelFileName, "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			fmt.Print("Total Lines ", sheet1.MaxRow, sheet1.Name)
			for i := 7; i <= (int(sheet1.MaxRow)); i++ {
				contar++
				row1 := sheet1.Row(i)
				col1 := row1.Col(0)
				if contar > 1 && strings.Trim(col1, " ") != "Totales Bs.:" {
					coma = ","
				} else {
					coma = ""
				} //Fin de contar para colocar coma
				re := regexp.MustCompile(`[-()]`)
				agen := re.Split(row1.Col(2), -1)
				agencia, venta := strings.ToUpper(agen[0]), strings.Trim(row1.Col(4), " ")
				premio, comision := strings.Trim(row1.Col(6), " "), strings.Trim(row1.Col(5), " ")
				insertar += coma
				insertar += "('" + agencia + "'," + venta + "," + premio + "," + comision
				insertar += ",1,'" + a.Fecha + "',Now()," + strconv.Itoa(posicionarchivo) + "," + strconv.Itoa(oid) + ")"
				a.Salvar = true
			} //Fin de Repita para
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

func (a *Archivo) LeerMaticloXLSX(ch chan []byte, tipo string) (bool, string) {

	fig := SLoteria
	posicionarchivo := 5
	if tipo == "f" {
		fig = SFigura
		posicionarchivo = 13
	}
	a.iniciarVariable(fig)

	insertar := a.Cabecera
	var coma string
	contar := 0
	oid, b := a.CrearTraza(posicionarchivo, a.ConvertirTablaNumero(fig))
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
			if a.CantidadLineas > 7 {
				contar++
				for _, cell := range row.Cells {
					text := cell.String()
					if strings.Trim(text, " ") != "" {
						cel = append(cel, text)
					}
				} //FIN DE LA CELDA

				l := len(cel)
				if l > 7 {
					if contar > 1 && strings.Trim(cel[0], " ") != "Totales Bs.:" {
						coma = ","
					} else {
						coma = ""
					}
					re := regexp.MustCompile(`[-()]`)
					agen := re.Split(cel[2], -1)
					agencia, venta := strings.ToUpper(agen[0]), strings.Trim(cel[4], " ")
					premio, comision := strings.Trim(cel[6], " "), strings.Trim(cel[5], " ")
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
