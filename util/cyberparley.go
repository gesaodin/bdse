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
func (a *Archivo) LeerCyberParley(ch chan []byte, tipo string) (bool, string) {
	posicionarchivo, fig := tipoCyberParley(tipo)
	a.iniciarVariable(fig)

	insertar := a.Cabecera
	var coma string

	oid, b := a.CrearTraza(posicionarchivo, a.ConvertirTablaNumero(fig))
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

		if strings.Trim(linea[0], " ") == "Fecha" {
			a.Leer = true
			a.CantidadLineas++
		}
		if a.Leer {
			if a.CantidadLineas > 2 && strings.Trim(linea[1], " ") == "Agencia" {
				coma = ","
			} else {
				coma = ""
			}
			insertar += coma
			if a.CantidadLineas > 1 && strings.Trim(linea[1], " ") == "Agencia" {

				re := regexp.MustCompile(`[(-)]`)
				agen := re.Split(linea[2], -1)
				// fmt.Println(agen[2])

				c := strings.Replace(strings.Trim(linea[5], " "), "-", "", -1)
				p := strings.Replace(strings.Trim(linea[6], " "), "-", "", -1)
				agencia, venta := strings.Trim(agen[0], " "), strings.Trim(linea[4], " ")
				premio, comision := p, c
				insertar += "('" + agencia + "'," + venta + "," + premio + ","
				insertar += comision + ",1,'" + a.Fecha + "',Now()," + strconv.Itoa(posicionarchivo) + "," + strconv.Itoa(oid) + ")"
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

func tipoCyberParley(tipo string) (posicionarchivo int, fig string) {
	switch tipo {
	case "p":
		posicionarchivo = 7
		fig = SParley
		break
	case "f":
		posicionarchivo = 25
		fig = SFigura
		break
	default:
		break
	}
	return
}
