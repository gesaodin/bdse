//utilidades generales para cadenas y números
package util

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gesaodin/bdse/util/logs/mensaje"
	"github.com/tealeg/xlsx"
)

const (
	//Loteria terminales y triples
	Loteria string = "0"
	//Parley apuestas generales deporte
	Parley string = "1"
	//Totales
	_Totales = "TOTALES"
)

//webMSJ Mensajes del sistema
type webMSJ struct {
	Tipo    int    `json:"tipo"`
	Mensaje string `json:"msj"`
	Autor   string `json:"aut"`
}

//Archivo Estructura de los archivos
type Archivo struct {
	Responsable      int
	Ruta             string
	NombreDelArchivo string
	Codificacion     string
	Cabecera         string
	Leer             bool
	Salvar           bool
	Fecha            string
	CantidadLineas   int
	Registros        int
	PostgreSQL       *sql.DB
	Canal            chan []byte
}

var m mensaje.WChat

//iniciarVariable Variables del sistema
func (a *Archivo) iniciarVariable(tabla string) {
	a.Cabecera = "INSERT INTO " + tabla + " (agen,vent,prem,comi,usua,fech,fcre,sist, arch) VALUES "
	a.CantidadLineas = 0
	a.Leer = false
	a.Salvar = false
}

//CrearTraza Traza y eventos de los archivos
func (a *Archivo) CrearTraza(tipo int, tabla string) (oid int, err error) {
	t := time.Now()
	nomb := a.NombreDelArchivo
	urls := a.Ruta
	resp := strconv.Itoa(a.Responsable)
	fcre := t.Format("2006-01-02 15:04:05")
	s := "INSERT INTO archivo (esta,nomb,fech, fcre,urls,resp,publ,tipo, tabl) VALUES "
	s += "(0,'" + nomb + "','" + a.Fecha + "','" + fcre + "','" + urls + "'," + resp + ",1,"
	s += strconv.Itoa(tipo) + "," + tabla + ") RETURNING oid"

	sq, err := a.PostgreSQL.Query(s)
	if err != nil {

		return 0, err
	}

	for sq.Next() {
		sq.Scan(&oid)
	}

	return
}

//ModificarTraza Actualizar Traza o eventos
func (a *Archivo) ModificarTraza() bool {
	t := time.Now()
	fpro := t.Format("2006-01-02 15:04:05")
	nomb := a.NombreDelArchivo
	s := "UPDATE archivo SET cant=" + strconv.Itoa(a.Registros) + ", esta=1,fpro='" + fpro + "' WHERE nomb='" + nomb + "'"
	//fmt.Println(s)
	_, err := a.PostgreSQL.Exec(s)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

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
	//fmt.Println(ch)
	a.Canal <- j
	//web.Mensajeria.Usuario["gpanel"].ch <- j

	return true, insertar
}

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

//LeerMaticlo Archivo en formato XLS 97-2003
func (a *Archivo) LeerMaticlo(ch chan []byte) (bool, string) {
	a.iniciarVariable("loteria")
	insertar := a.Cabecera
	var coma string
	contar := 0

	oid, b := a.CrearTraza(5, Loteria)
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
					insertar += coma
					re := regexp.MustCompile(`[-()]`)
					agen := re.Split(cel[2], -1)
					agencia, venta := strings.ToUpper(agen[0]), strings.Trim(cel[4], " ")
					premio, comision := strings.Trim(cel[6], " "), strings.Trim(cel[5], " ")
					insertar += "('" + agencia + "'," + venta + "," + premio + "," + comision
					insertar += ",1,'" + a.Fecha + "',Now(),5," + strconv.Itoa(oid) + ")"
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

//LeerIlbanquero Consultar en ilbanquero con el proveedor
func (a *Archivo) LeerIlbanquero(ch chan []byte) (bool, string) {
	a.iniciarVariable("parley")
	insertar := a.Cabecera
	var coma string

	oid, b := a.CrearTraza(6, Parley)
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
			if "Taquillas" == strings.Trim(linea[0], " ") {
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
					re := regexp.MustCompile(`[:]`)
					agen := re.Split(linea[0], -1)
					agenc := strings.Split(strings.Trim(agen[2], " "), " ")
					agencia, venta := strings.Trim(agenc[0], " "), strings.Trim(linea[1], " ")
					premio, comision := strings.Trim(linea[4], " "), strings.Trim(linea[6], " ")
					insertar += "('" + agencia + "'," + venta + "," + premio + "," + comision
					insertar += ",1,'" + a.Fecha + "',Now(),6," + strconv.Itoa(oid) + ")"
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

//LeerSport Sport
func (a *Archivo) LeerSport(ch chan []byte) (bool, string) {
	a.iniciarVariable("parley")
	insertar := a.Cabecera
	var coma string
	oid, b := a.CrearTraza(8, Parley)
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
					// fmt.Println(agen[2])

					p := strings.Replace(strings.Trim(linea[5], " "), "-", "", -1)
					c := strings.Replace(strings.Trim(linea[2], " "), "-", "", -1)
					agencia, venta := strings.Trim(agen[0], " "), strings.Trim(linea[1], " ")
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

//LeerTodo Todo un archivo
func (a *Archivo) LeerTodo() (f []byte, err error) {
	f, err = ioutil.ReadFile(a.NombreDelArchivo)
	return
}

//LeerCodigosYCrearAgencias Crear codigos y agencias
func (a *Archivo) LeerCodigosYCrearAgencias() bool {
	var sql string
	archivo, err := os.Open("public/temp/Com-Gru-Age-Caja.csv")
	Error(err)

	scan := bufio.NewScanner(archivo)
	for scan.Scan() {
		linea := strings.Split(scan.Text(), ";")

		grupo := linea[1]
		sql = `INSERT INTO grupo (comer, obse,tipo) VALUES (1,'` + grupo + `',0);`
		_, err := a.PostgreSQL.Exec(sql)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(sql)

		cap := linea[2]

		dondegrupo := `(SELECT oid FROM grupo WHERE obse='` + grupo + `')`
		dondeagencia := `(SELECT oid FROM agencia WHERE obse='` + cap + `')`

		sql = `INSERT INTO agencia (comer,grupo,subgr,colec,obse)
		VALUES (1,` + dondegrupo + `,0,0,'` + cap + `');`
		_, err = a.PostgreSQL.Exec(sql)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(sql)

		caja := linea[3]
		sql = `INSERT INTO zr_agencia (comer,grupo,subgr,colec,oida,codi)
		VALUES (1,` + dondegrupo + `,0,0,` + dondeagencia + `,'` + caja + `'); `
		_, err = a.PostgreSQL.Exec(sql)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(sql)
		sql = `INSERT INTO usuario (nomb,ncom,corr,fech,esta,rol, toke) VALUES
					(
						'` + grupo + `', 'Grupo','agencia@admin.com',
						Now(), 2, 'Grupo', md5('` + grupo + `123456')

					)`
		_, err = a.PostgreSQL.Exec(sql)
		if err != nil {
			fmt.Println(err.Error())
		}
		sql = `INSERT INTO usuario (nomb,ncom,corr,fech,esta,rol, toke) VALUES
					(
						'` + cap + `', 'Agencia','agencia@admin.com',
						Now(), 4, 'Agencia', md5('` + cap + `123456')

					)`
		_, err = a.PostgreSQL.Exec(sql)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return true
}

//LeerCodigosYCrearSaldos Saldos del sistema
func (a *Archivo) LeerCodigosYCrearSaldos() bool {
	var sql string
	archivo, err := os.Open("public/temp/saldos.enero.csv")
	Error(err)

	scan := bufio.NewScanner(archivo)
	for scan.Scan() {
		linea := strings.Split(scan.Text(), ";")

		cap := linea[0]
		saldo := linea[1]
		dondeagencia := `(SELECT oid FROM agencia WHERE obse='` + cap + `')`
		sql = `INSERT INTO cobrosypagos (oida, fech, vien) VALUES (` + dondeagencia + `,'2016-12-31'::TIMESTAMP,` + saldo + `);`
		_, err := a.PostgreSQL.Exec(sql)
		if err != nil {
			fmt.Println(err.Error())
		}

		//fmt.Println(sql)

	}
	sql = `INSERT INTO cobrosypagos_grupo (oidg,fech,vien,sald,movi,van,erec)
	select gr.oid, '2016-12-31'::TIMESTAMP, sum(cyp.vien),0,0,sum(cyp.vien),0 from grupo gr
	JOIN agencia ag ON ag.grupo=gr.oid
	JOIN cobrosypagos cyp ON ag.oid=cyp.oida
	GROUP BY gr.oid`
	_, err = a.PostgreSQL.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
	}
	return true
}

func (a *Archivo) LeerEntregados() bool {
	var sql string

	archivo, err := os.Open("public/temp/EOAG012017.csv")
	Error(err)
	var i int
	scan := bufio.NewScanner(archivo)
	for scan.Scan() {
		var id int
		i++
		linea := strings.Split(scan.Text(), ";")
		if i > 1 {
			fecha := linea[0]
			codcuentas := linea[3]
			cuenta := linea[6]
			agencia := linea[8]
			voucher := linea[7]
			observacion := linea[13] + "|" + cuenta + "|" + linea[9] + "|" + linea[10] + "|" + linea[11]
			monto := ConvertirMonedaANumero(linea[12])
			//fmt.Println(agencia)
			sql = `SELECT oid FROM agencia WHERE obse='` + agencia + `'`
			//fmt.Println(sql)
			row, err := a.PostgreSQL.Query(sql)
			if err != nil {
				fmt.Println("A ocurrido un error en la conexion")
			}
			for row.Next() {
				row.Scan(&id)
			}

			if id != 0 {
				s := `INSERT INTO debe (comer,grupo,subgr,colec,oida,agen,mont,vouc,banc,fdep,freg,fope,fapr,tipo,esta,obse) VALUES
								(1,0,0,0,` + strconv.Itoa(id) + `,'` + agencia + `',` + monto + `,'` + voucher + `',` +
					codcuentas + `,'` + fecha + `','` + fecha + `','` + fecha + `'::DATE -1, '` + fecha + `',1,1,'` + observacion + `');`
				_, err := a.PostgreSQL.Exec(s)
				if err != nil {
					fmt.Println(s)
					fmt.Println(err.Error())
				}

			}

		}

	}

	return true
}

func (a *Archivo) LeerEntregadosGrupo() bool {
	var sql string

	archivo, err := os.Open("public/temp/EOGR012017.csv")
	Error(err)
	var i int
	scan := bufio.NewScanner(archivo)
	for scan.Scan() {
		var id int
		i++
		linea := strings.Split(scan.Text(), ";")
		if i > 1 {
			fecha := linea[0]
			codcuentas := linea[3]
			cuenta := linea[6]
			grupo := linea[8]
			voucher := linea[7]
			observacion := linea[13] + "|" + cuenta + "|" + linea[9] + "|" + linea[10] + "|" + linea[11]
			monto := ConvertirMonedaANumero(linea[12])
			fmt.Println(grupo)
			sql = `SELECT oid FROM grupo WHERE obse='` + grupo + `'`
			row, err := a.PostgreSQL.Query(sql)
			if err != nil {
				fmt.Println("A ocurrido un error en la conexion")
			}
			for row.Next() {
				row.Scan(&id)
			}

			if id != 0 {
				s := `INSERT INTO debe (comer,grupo,subgr,colec,oida,agen,mont,vouc,banc,fdep,freg,fope,fapr,tipo,esta,obse) VALUES
								(1,` + strconv.Itoa(id) + `,0,0,0,'',` + monto + `,'` + voucher + `',` +
					codcuentas + `,'` + fecha + `','` + fecha + `','` + fecha + `'::DATE -1, '` + fecha + `',1,1,'` + observacion + `');`
				_, err := a.PostgreSQL.Exec(s)
				if err != nil {
					fmt.Println(s)
					fmt.Println(err.Error())
				}

			}

		}

	}

	return true
}

func (a *Archivo) LeerEntregadosOficina() bool {

	archivo, err := os.Open("public/temp/EOAD012017.csv")
	Error(err)
	var i int
	scan := bufio.NewScanner(archivo)
	for scan.Scan() {
		i++
		linea := strings.Split(scan.Text(), ";")

		if i > 1 {
			var mov Movimiento
			mov.Fecha = linea[0]

			mov.CuentaDebe, _ = strconv.Atoi(linea[3])
			mov.TipoDebe = 1

			mov.CuentaHaber, _ = strconv.Atoi(linea[8])
			mov.TipoHaber = 1
			mov.Voucher = linea[7]

			mov.Monto = ConvertirMonedaANumero(linea[12])
			mov.Observacion = linea[14] + "|" + linea[10] + "|" + linea[11]

			ingreso, egreso := mov.generarSQL()
			//fmt.Println(sql)
			_, err = a.PostgreSQL.Exec(ingreso + egreso)

			if err != nil {
				fmt.Println(err.Error())
			}
		}

	}

	return true
}

type Movimiento struct {
	Oid               int     `json:"oid,omitempty"`
	Comercializadora  int     `json:"comercializadora,omitempty"`
	Grupo             int     `json:"grupo,omitempty"`
	SubGrupo          int     `json:"subgrupo,omitempty"`
	Colector          int     `json:"colector,omitempty"`
	AgenciaCod        int     `json:"agenciacod,omitempty"`
	Agencia           string  `json:"agencia,omitempty"`
	Nombre            string  `json:"nombre,omitempty"`
	Fecha             string  `json:"fecha,omitempty"`
	FDeposito         string  `json:"fdeposito,omitempty"`
	FOperacion        string  `json:"foperacion,omitempty"`
	Voucher           string  `json:"voucher,omitempty"`
	FormaDePago       int     `json:"forma,omitempty"`
	TipoDeOperacion   int     `json:"operacion,omitempty"`
	TipoTabla         int     `json:"tipo,omitempty"`
	Monto             string  `json:"monto,omitempty"`
	Cuota             float64 `json:"cuota,omitempty"`
	Cuenta            int     `json:"cuenta,omitempty"`
	CuentaDebe        int     `json:"cuentadebe,omitempty"`
	CuentaDebeNombre  string  `json:"cuentadeben,omitempty"`
	TipoDebe          int     `json:"tipodebe,omitempty"`
	CuentaHaber       int     `json:"cuentahaber,omitempty"`
	CuentaHaberNombre string  `json:"cuentahabern,omitempty"`
	TipoHaber         int     `json:"tipohaber,omitempty"`
	Banco             int     `json:"banco,omitempty"`
	BancoNombre       string  `json:"banconombre,omitempty"`
	Estatus           int     `json:"estatus,omitempty"`
	Observacion       string  `json:"observacion,omitempty"`
	Token             string  `json:"token,omitempty"`
}

//generarSQL Consultar
func (m *Movimiento) generarSQL() (sqlI string, sqlE string) {
	sql1 := "INSERT INTO "
	ie := "(comer,grupo,subgr,colec,agenc,fech,freg,tipo,cuen,mont,oper,obse, toke)" // INGRESO | EGRESO

	iii := "(" + strconv.Itoa(m.Comercializadora) + "," + strconv.Itoa(m.Grupo)
	iii += "," + strconv.Itoa(m.SubGrupo) + "," + strconv.Itoa(m.Colector) + "," + strconv.Itoa(m.AgenciaCod)
	iii += ",'" + m.Fecha + "',now(),"
	cuenta := strconv.Itoa(m.TipoDebe) + "," + strconv.Itoa(m.CuentaDebe) + ","
	iff := m.Monto + ", '" + m.Voucher + "', '" + m.Observacion + "', md5('" + m.Fecha + m.Voucher + m.Monto + "'));"
	sqlI = sql1 + "movimiento_ingreso " + ie + " VALUES " + iii + cuenta + iff

	cuenta = strconv.Itoa(m.TipoHaber) + "," + strconv.Itoa(m.CuentaHaber) + ","
	sqlE = sql1 + "movimiento_egreso " + ie + " VALUES " + iii + cuenta + iff
	//sqlE = sqls + sqle

	return
}
