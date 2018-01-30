//utilidades generales para cadenas y números
//Los archivos de Pos, Aliens y Turco son el mismo formato
package util

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gesaodin/bdse/util/logs/mensaje"
)

const (
	//Loteria terminales y triples
	Loteria string = "0"
	//Parley apuestas generales deporte
	Parley string = "1"
	//Animalitos apuestas generales animales
	Animalitos string = "2"
	//Figuras apuestas generales animales
	Figuras string = "3"
	//Caballos apuestas generales animales
	Caballos string = "4"
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

		_, err = a.PostgreSQL.Exec(sql)
		if err != nil {
			fmt.Println(err.Error())
		}

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

//LeerEntregados Saldos del sistema
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

//LeerEntregadosGrupo Saldos del sistema
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

//LeerEntregadosOficina Saldos del sistema
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
