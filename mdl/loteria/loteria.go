//juegos y programas para el azar en triples y terminales
package loteria

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/util"
)

//Archivo formatos en estructuras
type Archivo struct {
	Oid       int       `json:"oid,omitempty"`
	Nombre    string    `json:"nombre,omitempty"`
	Fecha     time.Time `json:"fecha,omitempty"`
	Creado    time.Time `json:"creado,omitempty"`
	Procesado time.Time `json:"procesado,omitempty"`
	Estatus   string    `json:"estatus,omitempty"`
	Tabla     string    `json:"tabla,omitempty"`
	Cantidad  int       `json:"cantidad,omitempty"`
	idTabla   int
}

//Reporte reglas para la impresión
type Reporte struct {
	Agencia  string  `json:"age,omitempty"`
	Venta    float32 `json:"ven,omitempty"`
	Premio   float32 `json:"pre,omitempty"`
	Comision float32 `json:"com,omitempty"`
	Saldo    float32 `json:"sal,omitempty"`
	Fecha    string  `json:"fec,omitempty"`
	Tabla    int     `json:"tabla,omitempty"`
	Archivo  int     `json:"archivo,omitempty"`
	Sistema  int     `json:"sistema,omitempty"`
	Usuario  int     `json:"usuario,omitempty"`
}

//Listar valores generales
type Listar struct {
	Oid    int    `json:"oid,omitempty"`
	Nombre string `json:"nombre,omitempty"`
	Tipo   int    `json:"tipo,omitempty"`
}

//JsonDataReporte esquemas de impresión
type JsonDataReporte struct {
	Id      int    `json:"id,omitempty"`
	Tabla   string `json:"tabla,omitempty"`
	Desde   string `json:"desde,omitempty"`
	Hasta   string `json:"hasta,omitempty"`
	Sistema int    `json:"sistema,omitempty"`
	Tipo    int    `json:"tipo,omitempty"`
}

//ArchivosCargados listar todos los registros de los archivos
func (r *Reporte) ArchivosCargados(data JsonDataReporte) (j []byte, e error) {
	var lst []Archivo

	s := "SELECT oid,esta,nomb,fech, fcre,fpro, tabl, cant FROM archivo WHERE fech BETWEEN "
	s += "'" + data.Desde + " 00:00:00'::TIMESTAMP AND '" + data.Hasta + " 23:59:59'::TIMESTAMP"

	row, e := sys.PostgreSQL.Query(s)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	defer row.Close()
	for row.Next() {
		var a Archivo
		var oid int
		var tabl, cant int

		var esta, nomb sql.NullString
		var fech, fcre, fpro time.Time

		row.Scan(&oid, &esta, &nomb, &fech, &fcre, &fpro, &tabl, &cant)
		a.Oid = oid
		a.Estatus = util.ValidarNullString(esta)
		a.Nombre = util.ValidarNullString(nomb)
		a.Fecha = fech
		a.Creado = fcre
		a.idTabla = tabl
		a.ConvertirTabla()
		a.Cantidad = cant
		a.Procesado = fpro
		//a.Procesado = util.ValidarNullString(fpro)
		lst = append(lst, a)

	}
	j, _ = json.Marshal(lst)

	return
}

//Saldos ver los montos acumulados
func (r *Reporte) Saldos(data JsonDataReporte) (j []byte, e error) {
	var lst []interface{}
	var s, donde string

	tbl := "loteria"
	if data.Tabla != "" {
		tbl = data.Tabla
	}
	fmt.Println(data.Tabla)
	if data.Sistema > 0 && data.Sistema != 99 {
		donde = " AND sist = " + strconv.Itoa(data.Sistema)
	}
	s = "SELECT agen, fech, vent-prem-comi as saldo, vent, prem, comi "
	s += "FROM " + tbl + " WHERE arch = " + strconv.Itoa(data.Id) + donde
	fmt.Println(s)
	if data.Id == 0 {
		s = "SELECT agen, fech, vent-prem-comi as saldo, vent, prem, comi FROM " + tbl + " WHERE fech BETWEEN "
		s += "'" + data.Desde + " 00:00:00'::TIMESTAMP AND '" + data.Hasta + " 23:59:59'::TIMESTAMP"
		s += donde
		fmt.Println(s)
	}

	row, e := sys.PostgreSQL.Query(s)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	defer row.Close()
	for row.Next() {
		var rp Reporte
		var saldo float32
		var fech time.Time
		var agen string

		var vent, prem, comi float32
		row.Scan(&agen, &fech, &saldo, &vent, &prem, &comi)
		rp.Venta = vent
		rp.Premio = prem
		rp.Comision = comi

		rp.Agencia = agen
		rp.Fecha = fech.String()[0:10]
		rp.Saldo = saldo

		lst = append(lst, rp)
	}
	// println(lst)
	j, _ = json.Marshal(lst)
	return
}

//Sistemas listar tabla de programas
func (l *Listar) Sistemas(data JsonDataReporte) (j []byte, e error) {
	var lst []Listar
	var donde string
	if data.Id != 3 {
		donde = " WHERE arch =" + strconv.Itoa(data.Id)
	}
	s := "SELECT oid, obse, arch FROM sistema " + donde + " ORDER BY arch, oid"
	row, e := sys.PostgreSQL.Query(s)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	defer row.Close()
	for row.Next() {
		var ls Listar
		var oid, arch int

		var nomb string
		row.Scan(&oid, &nomb, &arch)
		ls.Oid = oid
		ls.Nombre = nomb
		ls.Tipo = arch
		lst = append(lst, ls)
	}
	j, _ = json.Marshal(lst)
	return
}

//Banca comercializadora
type Banca struct {
	Reporte []Reporte
}

//SaldosGenerales consultar el saldo total acumulado
func (l *Listar) SaldosGenerales(data JsonDataReporte) (jSon []byte, e error) {
	var banca map[string]Banca
	var nombreAuxiliar string
	// var agenciaAuxiliar string
	banca = make(map[string]Banca)
	s := `
				SELECT agencia.obse AS nombre, t.agen AS agencia,
				t.vent AS venta, t.prem AS premio, t.comi AS comision,
				t.saldo AS saldo,  archivo.tabl AS tabla,
				t.arch AS archivo, t.sist AS sistema, t.fech AS fecha FROM
					(
						select arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria UNION
						select arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
					) AS t
				JOIN archivo ON archivo.oid=t.arch
				JOIN sistema ON sistema.oid=archivo.tipo
				LEFT JOIN zr_agencia ON t.agen=zr_agencia.codi
				LEFT JOIN agencia ON agencia.oid=zr_agencia.oida
				ORDER BY agencia.oid, t.sist`

	row, e := sys.PostgreSQL.Query(s)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	defer row.Close()
	var i int
	var reporte []Reporte
	for row.Next() {

		var nombre, agencia sql.NullString
		var venta, premio, comision, saldo float32
		var tabla, archivo, sistema int
		var fecha time.Time
		var r Reporte

		row.Scan(&nombre, &agencia, &venta, &premio, &comision, &saldo, &tabla, &archivo, &sistema, &fecha)

		r.Agencia = util.ValidarNullString(agencia)
		r.Venta = venta
		r.Premio = premio
		r.Comision = comision
		r.Saldo = saldo
		r.Sistema = sistema
		r.Fecha = fecha.String()[0:10]

		if i == 0 {
			nombreAuxiliar = util.ValidarNullString(nombre)
			// agenciaAuxiliar = util.ValidarNullString(agencia)

		}

		//fmt.Println(nombreAuxiliar, nombre)
		if nombreAuxiliar != util.ValidarNullString(nombre) {
			var b Banca
			b.Reporte = reporte
			banca[nombreAuxiliar] = b
			nombreAuxiliar = util.ValidarNullString(nombre)
			reporte = nil
		}

		// if agenciaAuxiliar != util.ValidarNullString(agencia) {
		// 	agenciaAuxiliar = util.ValidarNullString(agencia)
		//
		// }

		reporte = append(reporte, r)

		i++
	}
	var b Banca
	b.Reporte = reporte
	banca[nombreAuxiliar] = b
	//
	//fmt.Printf("%v", banca)
	jSon, _ = json.Marshal(banca)

	return
}

//Detalle establece las reglas generales
type Detalle struct {
	Sistema int     `json:"sistema,omitempty"`
	Tabla   int     `json:"tabla,omitempty"`
	Saldo   float32 `json:"saldo"`
	Fecha   string  `json:"fecha,omitempty"`
}

//SaldosGeneralesPorSistemas establece los montos acumulados por los programas
func (r *Reporte) SaldosGeneralesPorSistemas(data JsonDataReporte) (jSon []byte, e error) {
	var lst map[string]interface{}
	lst = make(map[string]interface{})
	var fechaAuxiliar string
	var dtl []Detalle

	s := `
	SELECT sistema.oid, t.fech,  sum(t.saldo), tabl FROM sistema
		JOIN (
			select arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria
		UNION
			select arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
		) AS t ON sistema.oid=t.sist
	JOIN archivo ON archivo.oid=t.arch
	WHERE `
	s += "t.fech BETWEEN "
	s += "'" + data.Desde + " 00:00:00'::TIMESTAMP AND '" + data.Hasta + " 23:59:59'::TIMESTAMP"
	s += " AND tabl=" + strconv.Itoa(data.Id)

	s += " GROUP BY sistema.oid, t.fech, tabl ORDER BY t.fech, sistema.oid"

	row, e := sys.PostgreSQL.Query(s)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	defer row.Close()
	var i int

	for row.Next() {
		var oid, tabla int
		var saldo float32
		var fech sql.NullString

		var d Detalle
		row.Scan(&oid, &fech, &saldo, &tabla)
		d.Sistema = oid
		d.Saldo = saldo
		d.Tabla = tabla

		if i == 0 {
			fechaAuxiliar = util.ValidarNullString(fech)[0:10]
		}
		if fechaAuxiliar != util.ValidarNullString(fech)[0:10] {
			lst[fechaAuxiliar] = dtl
			dtl = nil
			fechaAuxiliar = util.ValidarNullString(fech)[0:10]
		}
		dtl = append(dtl, d)
		i++
	}

	lst[fechaAuxiliar] = dtl

	jSon, _ = json.Marshal(lst)
	return
}

//SaldoGeneralTotales totales acumulados por fecha
func (r *Reporte) SaldoGeneralTotales(data JsonDataReporte) (jSon []byte, e error) {
	var i int
	var tablaAuxiliar int
	var dtl []Detalle
	var lst map[int]interface{}
	lst = make(map[int]interface{})

	s := `SELECT t.fech,  sum(t.saldo) AS saldo, tabl AS tabla FROM sistema
		JOIN (
			select arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria
		UNION
			select arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
		) AS t ON sistema.oid=t.sist
	JOIN archivo ON archivo.oid=t.arch
	WHERE `
	s += " t.fech BETWEEN "
	s += "'" + data.Desde + " 00:00:00'::TIMESTAMP AND '" + data.Hasta + " 23:59:59'::TIMESTAMP"
	s += " GROUP BY t.fech, tabl ORDER BY tabl, t.fech"

	row, e := sys.PostgreSQL.Query(s)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	defer row.Close()

	for row.Next() {
		var tabla int
		var saldo float32
		var fech sql.NullString
		var d Detalle
		row.Scan(&fech, &saldo, &tabla)

		d.Saldo = saldo
		d.Fecha = util.ValidarNullString(fech)[0:10]
		if i == 0 {
			tablaAuxiliar = tabla
		}

		if tablaAuxiliar != tabla {
			lst[tablaAuxiliar] = dtl
			dtl = nil
			tablaAuxiliar = tabla
		}
		dtl = append(dtl, d)
		i++
	}

	lst[tablaAuxiliar] = dtl
	jSon, _ = json.Marshal(lst)
	return
}

//Balance configuracion de saldos
type Balance struct {
	Saldo float32 `json:"saldo"`
	Debe  float64 `json:"debe"`
	Haber float64 `json:"haber"`
}

//BalanceGeneral reglas generales de los saldos
func (r *Reporte) BalanceGeneral(data JsonDataReporte) (jSon []byte, err error) {
	var lstBalance map[string]interface{}
	var lstTabla map[int]interface{}
	// var Detalle []interface{}
	var tablaAuxiliar int
	// var fechaAuxiliar string
	lstBalance = make(map[string]interface{})
	lstTabla = make(map[int]interface{})

	s := `	SELECT f.fech AS fecha, f.tabl AS tabla, f.saldo, saldodebe.totaldebe AS debe, saldohaber.totalhaber AS haber FROM (
		SELECT t.fech,  sum(t.saldo) AS saldo, tabl FROM sistema
			JOIN (
				select arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria
			UNION
				select arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
			) AS t ON sistema.oid=t.sist
		JOIN archivo ON archivo.oid=t.arch
		WHERE `
	s += " t.fech BETWEEN "
	s += "'" + data.Desde + " 00:00:00'::TIMESTAMP AND '" + data.Hasta + " 23:59:59'::TIMESTAMP"
	s += `
			GROUP BY t.fech, tabl
			ORDER BY tabl, t.fech ) AS f
			LEFT JOIN ( SELECT fech, SUM(mont) AS totaldebe FROM debe GROUP BY fech) saldodebe ON saldodebe.fech=f.fech
			LEFT JOIN ( SELECT fech, SUM(mont) AS totalhaber FROM haber GROUP BY fech) saldohaber ON saldohaber.fech=f.fech`
	//fmt.Println(s)
	row, e := sys.PostgreSQL.Query(s)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	// var balanceAnterior Balance
	defer row.Close()
	i := 0
	for row.Next() {
		var tabla int
		var fecha sql.NullString
		var saldo float32
		var debe sql.NullFloat64
		var haber sql.NullFloat64
		err = row.Scan(&fecha, &tabla, &saldo, &debe, &haber)
		if err != nil {
			fmt.Println(err.Error())
		}

		debf64 := util.ValidarNullFloat64(debe)
		habf64 := util.ValidarNullFloat64(haber)

		var balance Balance
		balance = Balance{Saldo: saldo, Debe: debf64, Haber: habf64}

		if i == 0 {
			// balanceAnterior = Balance{Saldo: saldo, Debe: debe, Haber: haber}
			tablaAuxiliar = tabla
			//fechaAuxiliar = util.ValidarNullString(fecha)[0:10]
		}

		/*
			if fechaAuxiliar != util.ValidarNullString(fecha)[0:10] {
				lstBalance[fechaAuxiliar] = balanceAnterior //Detalle
				balanceAnterior = balance
				fechaAuxiliar = util.ValidarNullString(fecha)[0:10]

			}
		*/

		if tablaAuxiliar != tabla {

			lstTabla[tablaAuxiliar] = lstBalance
			lstBalance = nil
			lstBalance = make(map[string]interface{})
			tablaAuxiliar = tabla

		}

		lstBalance[util.ValidarNullString(fecha)[0:10]] = balance //Detalle
		//balanceAnterior = balance
		//fechaAuxiliar = util.ValidarNullString(fecha)[0:10]

		//Detalle = append(Detalle, balance)

		i++

	}
	//lstBalance[fechaAuxiliar] = Detalle
	lstTabla[tablaAuxiliar] = lstBalance

	jSon, _ = json.Marshal(lstTabla)
	//fmt.Println(s)
	return
}

//ConvertirTabla devuelve loteria / parley
func (a *Archivo) ConvertirTabla() {
	switch a.idTabla {
	case 0:
		a.Tabla = "loteria"
		break
	case 1:
		a.Tabla = "parley"
		break
	case 2:
		a.Tabla = "figura"
		break
	default:
		a.Tabla = "loteria"
		break
	}
}
