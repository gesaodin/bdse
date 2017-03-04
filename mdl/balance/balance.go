package balance

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/util"
)

type Pago struct {
	Agencia       string  `json:"agencia,omitempty"`
	Taquilla      string  `json:"taquilla,omitempty"`
	Voucher       string  `json:"voucher,omitempty"`
	Desde         string  `json:"desde,omitempty"`
	Hasta         string  `json:"hasta,omitempty"`
	Deposito      string  `json:"deposito,omitempty"`
	Banco         int     `json:"banco,omitempty"`
	BancoNombre   string  `json:"banconombre,omitempty"`
	Fecha         string  `json:"fecha,omitempty"`
	FechaAprobado string  `json:"fechaaprobado,omitempty"`
	FormaDePago   int     `json:"forma,omitempty"`
	Monto         float64 `json:"monto,omitempty"`
	Venta         float64 `json:"venta,omitempty"`
	Premio        float64 `json:"premio,omitempty"`
	Comision      float64 `json:"comision,omitempty"`
	Saldo         float64 `json:"saldo,omitempty"`
	Vienen        float64 `json:"vienen,omitempty"`
	Recibido      float64 `json:"recibido,omitempty"`
	Entregado     float64 `json:"entregado,omitempty"`
	Ingreso       float64 `json:"ingreso,omitempty"`
	Egreso        float64 `json:"egreso,omitempty"`
	Prestamo      float64 `json:"prestamo,omitempty"`
	Cuota         float64 `json:"cuota,omitempty"`
	Estatus       int     `json:"estatus,omitempty"`
	Observacion   string  `json:"observacion,omitempty"`
	Sistema       int     `json:"sistema,omitempty"`
	Archivo       int     `json:"archivo,omitempty"`
}

type Respuesta struct {
	Cantidad int64  `json:"cant"`
	Msj      string `json:"msj"`
}

func (p *Pago) Registrar(data Pago) (jSon []byte, err error) {
	monto := strconv.FormatFloat(data.Monto, 'f', 6, 64)
	tabla := "haber"
	if(data.FormaDePago == 0){
		tabla = "debe"
	}
	forma := strconv.Itoa(data.FormaDePago)
	banco := strconv.Itoa(data.Banco)
	estatus := strconv.Itoa(data.Estatus)
	aprobado := "";
	campo := ""
	if data.FechaAprobado != ""{
		campo = "fapr,"
		aprobado = "'" + data.FechaAprobado + "',"
	}
	s := "INSERT INTO " + tabla + " (agen,mont,vouc,fdep,freg," + campo + "tipo,banc,esta,obse) VALUES "
	s += "('" + data.Agencia + "'," + monto + ",'" + data.Voucher + "',"
	s += "'" + data.Deposito + "',now()," + aprobado + "" + forma
	s += "," + banco + "," + estatus + ",'" + data.Observacion + "');"

	rs, err := sys.PostgreSQL.Exec(s)
	if err != nil {
		return
	}
	var res Respuesta
	cantidad, _ := rs.RowsAffected()
	res.Cantidad = cantidad
	res.Msj = "Se inserto correctamente"
	jSon, _ = json.Marshal(res)

	return
}

func (p *Pago) ListarPagos(data Pago) (jSon []byte, err error) {
	var s string
	s = `
			SELECT fdep,vouc,fapr,esta,debe.mont,debe.resp, banco.nomb FROM agencia
			INNER JOIN debe ON debe.agen=agencia.obse
			INNER JOIN banco ON debe.banc=banco.oid
			WHERE agencia.obse='` + data.Agencia + `'
			ORDER BY debe.fdep`
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}


	var lst []interface{}
	for row.Next() {
		var fdep, vouc, fapr, resp, nomb sql.NullString
		var esta int
		var mont float64

		var pago Pago
		e := row.Scan(&fdep, &vouc, &fapr, &esta, &mont, &resp, &nomb)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		var dep = util.ValidarNullString(fdep)
		var apr = util.ValidarNullString(fapr)

		if(apr != "null"){
			apr = apr[0:10]
		}

		pago.Voucher = util.ValidarNullString(vouc)
		pago.Observacion = util.ValidarNullString(resp)
		pago.BancoNombre = util.ValidarNullString(nomb)
		pago.FechaAprobado = apr
		pago.Fecha = dep[0:10]
		pago.Monto = mont
		pago.Estatus = esta
		lst = append(lst, pago)
	}

	jSon, _ = json.Marshal(lst)
	return
}

func (p *Pago) GenerarCobrosYPagos(data Pago) (jSon []byte, err error) {
	var fecha string = time.Now().String()[0:10]
	var s string

	if data.Agencia != "" {
		s = generarCobrosYPagosAgencia(data)

	} else {
		s = generarCobrosYPagosGeneral(fecha)
		if data.Fecha != "" {
			fecha = data.Fecha
			s = generarCobrosYPagosGeneral(fecha)
		}
		fmt.Println(s)
	}
	row, err := sys.PostgreSQL.Query(s)

	if err != nil {
		return
	}
	var lst []interface{}

	for row.Next() {
		var agen string
		var saldo, entregado, recibido, ingreso, egreso, prestamo, vienen, cuota sql.NullFloat64
		e := row.Scan(&agen, &saldo, &entregado, &recibido, &ingreso, &egreso, &prestamo, &vienen, &cuota)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		var pago Pago
		pago.Agencia = agen
		if data.Agencia != "" {
			pago.Agencia = ""
			pago.Fecha = agen[0:10]
		}
		pago.Saldo = util.ValidarNullFloat64(saldo)
		pago.Ingreso = util.ValidarNullFloat64(ingreso)
		pago.Egreso = util.ValidarNullFloat64(egreso)
		pago.Prestamo = util.ValidarNullFloat64(prestamo)
		pago.Entregado = util.ValidarNullFloat64(entregado)
		pago.Recibido = util.ValidarNullFloat64(recibido)
		pago.Vienen = util.ValidarNullFloat64(vienen)
		pago.Cuota = util.ValidarNullFloat64(cuota)

		lst = append(lst, pago)
	}

	jSon, _ = json.Marshal(lst)

	return
}

//Generar Cobros y Pagos Administrador
func generarCobrosYPagosGeneral(fecha string) (s string) {

	if fecha != "" {
		fecha = ` lotepar.fech BETWEEN '` + fecha + ` 00:00:00'::TIMESTAMP AND '` + fecha + ` 23:59:59'::TIMESTAMP `
	}
	s = `
		SELECT saldo_agencia.obse, saldo_agencia.saldo,
				debe.monto AS entregado, haber.monto AS recibido,
				ingreso.monto AS ingreso, egreso.monto AS egreso,
				prestamo.monto AS prestamo,
				cobrosypagos.vien AS vienen,
				cobrosypagos.van AS van
			FROM (
			SELECT agencia.obse, lotepar.fech, SUM(lotepar.saldo) AS saldo
			FROM agencia
			JOIN zr_agencia ON agencia.oid=zr_agencia.oida
			JOIN (
				SELECT arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria
				UNION
				SELECT arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
			) AS lotepar ON zr_agencia.codi=lotepar.agen

			WHERE ` + fecha + `
			GROUP BY agencia.obse,lotepar.fech
			) saldo_agencia

			-- DEBE
			LEFT JOIN (
					SELECT agen, fapr, SUM(mont) AS monto FROM debe
					GROUP BY agen,fapr
			) AS debe ON
			debe.agen=saldo_agencia.obse AND debe.fapr=saldo_agencia.fech

			-- HABER
			LEFT JOIN (
				SELECT agen, fapr, SUM(mont) AS monto FROM haber
				GROUP BY agen,fapr
			) AS haber ON
			haber.agen=saldo_agencia.obse  AND haber.fapr=saldo_agencia.fech

			--INGRESO
			LEFT JOIN (
				SELECT agen, fech, SUM(mont) AS monto FROM movimiento_ingreso
				GROUP BY agen,fech
			)
			AS ingreso ON
			ingreso.agen=saldo_agencia.obse  AND ingreso.fech=saldo_agencia.fech

			-- EGRESO
			LEFT JOIN (
				SELECT agen, fech, SUM(mont) AS monto FROM movimiento_egreso
				GROUP BY agen,fech
			)
			AS egreso ON
			egreso.agen=saldo_agencia.obse  AND egreso.fech=saldo_agencia.fech

			-- PRESTAMOS
			LEFT JOIN (
				SELECT agen, fech, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
				GROUP BY agen,fech)
			AS prestamo ON
			prestamo.agen=saldo_agencia.obse AND prestamo.fech=saldo_agencia.fech

			-- VIENEN
			LEFT JOIN cobrosypagos ON cobrosypagos.fech=saldo_agencia.fech
			AND cobrosypagos.agen=saldo_agencia.obse
			ORDER BY saldo_agencia.fech
	`
	return
}

//Estado de Cuenta por agencia
func generarCobrosYPagosAgencia(data Pago) (s string) {
	var fecha string = ""
	if data.Desde != "" {
		fecha = ` AND lotepar.fech BETWEEN '` + data.Desde + ` 00:00:00'::TIMESTAMP AND '` + data.Hasta + ` 23:59:59'::TIMESTAMP`
	}
	s = `
		SELECT saldo_agencia.fech, saldo_agencia.saldo,
				debe.monto AS entregado, haber.monto AS recibido,
				ingreso.monto AS ingreso, egreso.monto AS egreso,
				prestamo.monto AS prestamo,
				cobrosypagos.vien AS vienen,
				cobrosypagos.van AS van
			FROM (
			SELECT agencia.obse, lotepar.fech, SUM(lotepar.saldo) AS saldo
			FROM agencia
			JOIN zr_agencia ON agencia.oid=zr_agencia.oida
			JOIN (
				SELECT arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria
				UNION
				SELECT arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
			) AS lotepar ON zr_agencia.codi=lotepar.agen

			WHERE agencia.obse='` + data.Agencia + `'` + fecha + `
			GROUP BY agencia.obse,lotepar.fech
			) saldo_agencia

			-- DEBE
			LEFT JOIN (
					SELECT agen, fapr, SUM(mont) AS monto FROM debe
					GROUP BY agen,fapr
			) AS debe ON
			debe.agen=saldo_agencia.obse AND debe.fapr=saldo_agencia.fech

			-- HABER
			LEFT JOIN (
				SELECT agen, fapr, SUM(mont) AS monto FROM haber
				GROUP BY agen,fapr
			) AS haber ON
			haber.agen=saldo_agencia.obse  AND haber.fapr=saldo_agencia.fech

			--INGRESO
			LEFT JOIN (
				SELECT agen, fech, SUM(mont) AS monto FROM movimiento_ingreso
				GROUP BY agen,fech
			)
			AS ingreso ON
			ingreso.agen=saldo_agencia.obse  AND ingreso.fech=saldo_agencia.fech

			-- EGRESO
			LEFT JOIN (
				SELECT agen, fech, SUM(mont) AS monto FROM movimiento_egreso
				GROUP BY agen,fech
			)
			AS egreso ON
			egreso.agen=saldo_agencia.obse  AND egreso.fech=saldo_agencia.fech

			-- PRESTAMOS
			LEFT JOIN (
				SELECT agen, fech, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
				GROUP BY agen,fech)
			AS prestamo ON
			prestamo.agen=saldo_agencia.obse AND prestamo.fech=saldo_agencia.fech

			-- VIENEN
			LEFT JOIN cobrosypagos ON cobrosypagos.fech=saldo_agencia.fech
			AND cobrosypagos.agen=saldo_agencia.obse
			ORDER BY saldo_agencia.fech
	`
	return
}

func (p *Pago) GenerarCobrosYPagosSistemas(data Pago) (jSon []byte, err error) {
	//fech, saldo,sist,agen,arch
	var s string
	s = generarCobrosYPagosSistemas(data)
	row, err := sys.PostgreSQL.Query(s)

	if err != nil {
		return
	}
	var lst []interface{}

	for row.Next() {
		var fech, agen sql.NullString
		var sist, arch int
		var saldo sql.NullFloat64
		e := row.Scan(&fech, &saldo, &sist, &agen, &arch)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		var pago Pago
		pago.Fecha = util.ValidarNullString(fech)[0:10]
		pago.Saldo = util.ValidarNullFloat64(saldo)
		pago.Sistema = sist
		pago.Observacion = util.ValidarNullString(agen)
		pago.Archivo = arch
		lst = append(lst, pago)
	}

	jSon, _ = json.Marshal(lst)

	return
}

//Estado de Cuenta por Loteria y Parley (Sistemas)
func generarCobrosYPagosSistemas(data Pago) (s string) {
	var fecha string = ""
	if data.Desde != "" {
		fecha = ` AND lotepar.fech BETWEEN '` + data.Desde + ` 00:00:00'::TIMESTAMP AND '` + data.Hasta + ` 23:59:59'::TIMESTAMP`
	}

	s = `
		SELECT
			lotepar.fech, SUM(lotepar.saldo) AS saldo, lotepar.sist,
			sistema.obse, sistema.arch
		FROM agencia
		JOIN zr_agencia ON agencia.oid=zr_agencia.oida
		JOIN (
		SELECT
			arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria
		UNION
		SELECT
			arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
		) AS lotepar ON zr_agencia.codi=lotepar.agen
		JOIN sistema ON lotepar.sist=sistema.oid
		WHERE agencia.obse='` + data.Agencia + `'` + fecha + `
		GROUP BY lotepar.sist, lotepar.fech, sistema.arch, sistema.obse
		ORDER BY sistema.arch
	`
	return
}


func (p *Pago) GenerarCobrosYPagosDetallados(data Pago) (jSon []byte, err error) {
	//fech, saldo,sist,agen,arch
	var s string
	s = generarCobrosYPagosDetallados(data)
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}
	var lst []interface{}

	for row.Next() {
		var taquilla, observacion sql.NullString
		var venta, premio, comision, saldo sql.NullFloat64
		e := row.Scan(&taquilla, &venta, &premio, &comision, &saldo, &observacion)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		var pago Pago
		pago.Taquilla = util.ValidarNullString(taquilla)
		pago.Venta = util.ValidarNullFloat64(venta)
		pago.Premio = util.ValidarNullFloat64(premio)
		pago.Comision = util.ValidarNullFloat64(comision)
		pago.Saldo = util.ValidarNullFloat64(saldo)
		pago.Observacion = util.ValidarNullString(observacion)
		lst = append(lst, pago)
	}
	jSon, _ = json.Marshal(lst)
	return
}

func generarCobrosYPagosDetallados(data Pago) (s string){
	var fecha string = ""
	var sistema string = "";
	if data.Desde != "" {
		fecha = ` AND lotepar.fech BETWEEN '` + data.Desde + ` 00:00:00'::TIMESTAMP AND '` + data.Hasta + ` 23:59:59'::TIMESTAMP`
	}
	if data.Sistema > 0 {
		sistema = ` AND sistema.arch = ` + strconv.Itoa(data.Sistema)
	}
	s = `
		SELECT
			lotepar.agen,  lotepar.vent AS venta,
			lotepar.prem AS premio, lotepar.comi AS comision,
			lotepar.saldo AS saldo, sistema.obse
		FROM agencia
		JOIN zr_agencia ON agencia.oid=zr_agencia.oida

		JOIN (
			SELECT
				arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria

			UNION
			SELECT
				arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley

			) AS lotepar ON zr_agencia.codi=lotepar.agen
		JOIN sistema ON lotepar.sist=sistema.oid
		WHERE agencia.obse='` + data.Agencia + `' ` + fecha + sistema + ` ORDER BY sistema.oid,lotepar.agen  `
		return
}
