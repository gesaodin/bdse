//cobros y pagos en movimiento y entregados/recibidos
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

//Pago Control de Pagos
type Pago struct {
	Oid            int     `json:"oid,omitempty"`
	Banca          int     `json:"banca,omitempty"`
	Grupo          int     `json:"grupo,omitempty"`
	SubGrupo       int     `json:"subgrupo,omitempty"`
	Colector       int     `json:"colector,omitempty"`
	Agencia        string  `json:"agencia,omitempty"`
	Taquilla       string  `json:"taquilla,omitempty"`
	Voucher        string  `json:"voucher,omitempty"`
	Desde          string  `json:"desde,omitempty"`
	Hasta          string  `json:"hasta,omitempty"`
	Deposito       string  `json:"deposito,omitempty"`
	Banco          int     `json:"banco,omitempty"`
	BancoNombre    string  `json:"banconombre,omitempty"`
	Fecha          string  `json:"fecha,omitempty"`
	FechaAprobado  string  `json:"fechaaprobado,omitempty"`
	FechaOperacion string  `json:"fechaoperacion,omitempty"`
	FormaDePago    int     `json:"forma,omitempty"`
	Monto          float64 `json:"monto,omitempty"`
	Venta          float64 `json:"venta,omitempty"`
	Premio         float64 `json:"premio,omitempty"`
	Comision       float64 `json:"comision,omitempty"`
	Saldo          float64 `json:"saldo,omitempty"`
	Vienen         float64 `json:"vienen,omitempty"`
	Van            float64 `json:"van,omitempty"`
	Recibido       float64 `json:"recibido,omitempty"`
	Entregado      float64 `json:"entregado,omitempty"`
	Ingreso        float64 `json:"ingreso,omitempty"`
	Egreso         float64 `json:"egreso,omitempty"`
	Prestamo       float64 `json:"prestamo,omitempty"`
	Cuota          float64 `json:"cuota,omitempty"`
	Estatus        int     `json:"estatus,omitempty"`
	Observacion    string  `json:"observacion,omitempty"`
	Sistema        int     `json:"sistema,omitempty"`
	Archivo        int     `json:"archivo,omitempty"`
	Cierre         int     `json:"cierre,omitempty"`
}

//Respuesta Generales
type Respuesta struct {
	Cantidad int64  `json:"cant"` // Cantidad de elementos
	Msj      string `json:"msj"`  // Mensaje almacenado
}

//Registrar Un pago por movimiento de ingreso o egreso
func (p *Pago) Registrar(data Pago) (jSon []byte, err error) {
	monto := strconv.FormatFloat(data.Monto, 'f', 6, 64)
	var agencia string
	tabla := "haber"
	if data.FormaDePago == 0 {
		tabla = "debe"
	}
	forma := strconv.Itoa(data.FormaDePago)
	banco := strconv.Itoa(data.Banco)
	estatus := strconv.Itoa(data.Estatus)
	aprobado := ""
	campo := ""
	campofechaperacopm := ""
	operacion := ""
	if data.FechaAprobado != "" {
		campo = "fapr,"
		aprobado = "'" + data.FechaAprobado + "',"
	}

	if data.FechaOperacion != "" {
		campofechaperacopm = "fope,"
		operacion = "'" + data.FechaOperacion + "',"
	}

	if data.Oid != 0 {
		agencia = "(SELECT obse FROM agencia WHERE oid=" + strconv.Itoa(data.Oid) + ")"
	}

	s := "INSERT INTO " + tabla + " (comer,grupo,subgr,colec,oida,agen,mont,vouc,fdep,freg," + campofechaperacopm + campo + "tipo,banc,esta,obse) VALUES "
	s += "(" + strconv.Itoa(data.Banca) + "," + strconv.Itoa(data.Grupo) + "," + strconv.Itoa(data.SubGrupo)
	s += "," + strconv.Itoa(data.Colector) + "," + strconv.Itoa(data.Oid)
	s += "," + agencia + "," + monto + ",'" + data.Voucher + "',"
	s += "'" + data.Deposito + "',now()," + operacion + aprobado + forma
	s += "," + banco + "," + estatus + ",'" + data.Observacion + "');"

	fmt.Println(s)
	rs, err := sys.PostgreSQL.Exec(s)
	if err != nil {
		fmt.Println(s)
		fmt.Println(err.Error())
		return
	}
	var res Respuesta
	cantidad, _ := rs.RowsAffected()
	res.Cantidad = cantidad
	res.Msj = "Se inserto correctamente"
	jSon, _ = json.Marshal(res)

	return
}

//ListarPagos Generales del sistema
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

		if apr != "null" {
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

//GenerarCobrosYPagos Generacion de Cobros y Pagos del sistema
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

	}
	//fmt.Println(s)
	row, err := sys.PostgreSQL.Query(s)

	if err != nil {
		fmt.Println(s)
		fmt.Println(err.Error())
		return
	}
	var lst []interface{}

	for row.Next() {
		var oid, esta int
		var agen string
		var saldo, entregado, recibido, ingreso, egreso, prestamo, vienen, van, cuota sql.NullFloat64
		e := row.Scan(&oid, &agen, &vienen, &saldo, &entregado, &recibido, &ingreso, &egreso, &prestamo, &cuota, &van, &esta)
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
		pago.Oid = oid
		pago.Saldo = util.ValidarNullFloat64(saldo)
		pago.Ingreso = util.ValidarNullFloat64(ingreso)
		pago.Egreso = util.ValidarNullFloat64(egreso)
		pago.Prestamo = util.ValidarNullFloat64(prestamo)
		pago.Entregado = util.ValidarNullFloat64(entregado)
		pago.Recibido = util.ValidarNullFloat64(recibido)
		pago.Vienen = util.ValidarNullFloat64(vienen)
		pago.Van = util.ValidarNullFloat64(van)
		pago.Cuota = util.ValidarNullFloat64(cuota)
		pago.Estatus = esta

		lst = append(lst, pago)
	}

	jSon, _ = json.Marshal(lst)

	return
}

//GenerarCierreDiario Generar Cierre Diario de las operaciones contables
func (p *Pago) GenerarCierreDiario(data Pago) (jSon []byte, err error) {

	var s, tabla string
	var r Respuesta

	s = generarCobrosYPagosGeneralCierre(data.Fecha)
	fmt.Println(s)
	_, err = sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}
	data.Cierre = 1
	data.GenerarCobrosYPagosGrupo() //Crear el cierre

	tabla = "cobrosypagoscierre"
	fecha := `'` + data.Fecha + ` 00:00:00'::TIMESTAMP + '1 day'`
	s = `INSERT INTO ` + tabla + ` (fech,esta) VALUES (` + fecha + `, 0);
				UPDATE ` + tabla + ` SET esta=1 WHERE fech='` + data.Fecha + ` 00:00:00'::TIMESTAMP`
	_, err = sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}
	tabla = "cobrosypagoscierre_grupo"
	s = `INSERT INTO ` + tabla + ` (fech,esta) VALUES (` + fecha + `, 0);
				UPDATE ` + tabla + ` SET esta=1 WHERE fech='` + data.Fecha + ` 00:00:00'::TIMESTAMP`
	_, err = sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}
	r.Msj = "Proceso exitoso"
	//Calcular Participacion
	s = gCPGrupoParticipacionDiaria(data.Fecha)
	_, err = sys.PostgreSQL.Query(s)
	if err != nil {
		r.Msj += ", no se encontrarón participaciones diarias."
		fmt.Println(s)
		fmt.Println(err.Error())
		//return
	}

	r.Msj = "Proceso exitoso"
	jSon, _ = json.Marshal(r)

	return
}

//generarCobrosYPagosGeneral Generar Cobros y Pagos Administrador
func generarCobrosYPagosGeneral(fecha string) (s string) {
	fechaAux := `'` + fecha + ` 00:00:00'::TIMESTAMP `
	if fecha != "" {
		fecha = ` lotepar.fech BETWEEN '` + fecha + ` 00:00:00'::TIMESTAMP AND '` + fecha + ` 23:59:59'::TIMESTAMP `
	}
	s = `
	SELECT oid, obse,
		vienen, saldo, entregado, recibido,ingreso, egreso, prestamo,cuota,
		saldo + vienen + (entregado - recibido) + (egreso - (ingreso+prestamo)) AS van, esta FROM (
			SELECT z.oid, z.obse,
			COALESCE(x.saldo,0) AS saldo,
			COALESCE(debe.monto,0) AS entregado,
			COALESCE(haber.monto,0) AS recibido,
			COALESCE(ingreso.monto,0) AS ingreso,
			COALESCE(egreso.monto,0) AS egreso,
			COALESCE(prestamo.monto,0) AS prestamo,
			COALESCE(prestamo.cuota,0) AS cuota,
			COALESCE(cobrosypagos.vien,0) AS vienen,
			COALESCE(cobrosypagoscierre.esta,0) AS esta
			FROM grupo g
			LEFT JOIN agencia z ON g.oid=z.grupo
			LEFT JOIN (
			SELECT agencia.oid, lotepar.fech, SUM(lotepar.saldo) AS saldo
			FROM agencia
			LEFT JOIN zr_agencia ON agencia.oid=zr_agencia.oida
			LEFT JOIN (
			SELECT agen, fech, vent-prem-comi as saldo from loteria
			UNION
			SELECT agen, fech, vent-prem-comi as saldo from parley
			) AS lotepar ON zr_agencia.codi=lotepar.agen

			WHERE ` + fecha + `
			GROUP BY agencia.oid,lotepar.fech ) AS x ON x.oid=z.oid

			-- DEBE
			LEFT JOIN (
			SELECT agen, fapr, SUM(mont) AS monto FROM debe
			GROUP BY agen,fapr
			) AS debe ON
			debe.agen=z.obse AND debe.fapr=` + fechaAux + `

			-- HABER
			LEFT JOIN (
			SELECT agen, fapr, SUM(mont) AS monto FROM haber
			GROUP BY agen,fapr
			) AS haber ON
			haber.agen=z.obse  AND haber.fapr=` + fechaAux + `

			--INGRESO
			LEFT JOIN (
			SELECT agen, fapr, SUM(mont) AS monto FROM movimiento_ingreso
			GROUP BY agen,fapr
			)
			AS ingreso ON
			ingreso.agen=z.obse  AND ingreso.fapr=` + fechaAux + `

			-- EGRESO
			LEFT JOIN (
			SELECT agen, fapr, SUM(mont) AS monto FROM movimiento_egreso
			GROUP BY agen,fapr
			)
			AS egreso ON
			egreso.agen=z.obse  AND egreso.fapr=` + fechaAux + `

			-- PRESTAMOS
			LEFT JOIN (
			SELECT agen, fech, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
			GROUP BY agen,fech)
			AS prestamo ON
			prestamo.agen=z.obse AND prestamo.fech=` + fechaAux + `

			-- VIENEN
			LEFT JOIN cobrosypagos ON cobrosypagos.fech=` + fechaAux + `
			AND cobrosypagos.oida=z.oid

			-- CIERRE
			LEFT JOIN cobrosypagoscierre ON cobrosypagoscierre.fech=` + fechaAux + `
			WHERE  g.obse='AGE. DIRECTAS'
			ORDER BY z.obse
			) AS A
	`
	return
}

//generarCobrosYPagosGeneralCierre Generar Cierre Diario seleccionar el día y le suma al siguiente
func generarCobrosYPagosGeneralCierre(fecha string) (s string) {
	fechaAux := `'` + fecha + ` 00:00:00'::TIMESTAMP `
	sumar := fecha
	if fecha != "" {
		fecha = ` lotepar.fech BETWEEN '` + fecha + ` 00:00:00'::TIMESTAMP AND '` + fecha + ` 23:59:59'::TIMESTAMP `
	}
	s = `INSERT INTO cobrosypagos (oida, fech, vien)
	SELECT oid, '` + sumar + `'::TIMESTAMP + '1 day' AS fech,
		saldo + vienen + (entregado - recibido) + (egreso - (ingreso+prestamo)) AS van FROM (
			SELECT z.oid, obse,
			COALESCE(x.saldo,0) AS saldo,
			COALESCE(debe.monto,0) AS entregado,
			COALESCE(haber.monto,0) AS recibido,
			COALESCE(ingreso.monto,0) AS ingreso,
			COALESCE(egreso.monto,0) AS egreso,
			COALESCE(prestamo.monto,0) AS prestamo,
			COALESCE(prestamo.cuota,0) AS cuota,
			COALESCE(cobrosypagos.vien,0) AS vienen
			FROM agencia AS z
			LEFT JOIN (
			SELECT agencia.oid, lotepar.fech, SUM(lotepar.saldo) AS saldo
			FROM agencia
			LEFT JOIN zr_agencia ON agencia.oid=zr_agencia.oida
			LEFT JOIN (
			SELECT agen, fech, vent-prem-comi as saldo from loteria
			UNION
			SELECT agen, fech, vent-prem-comi as saldo from parley
			) AS lotepar ON zr_agencia.codi=lotepar.agen

			WHERE ` + fecha + `
			GROUP BY agencia.oid,lotepar.fech ) AS x ON x.oid=z.oid

			-- DEBE
			LEFT JOIN (
			SELECT agen, fapr, SUM(mont) AS monto FROM debe
			GROUP BY agen,fapr
			) AS debe ON
			debe.agen=z.obse AND debe.fapr=` + fechaAux + `

			-- HABER
			LEFT JOIN (
			SELECT agen, fapr, SUM(mont) AS monto FROM haber
			GROUP BY agen,fapr
			) AS haber ON
			haber.agen=z.obse  AND haber.fapr=` + fechaAux + `

			--INGRESO
			LEFT JOIN (
			SELECT agen, fech, SUM(mont) AS monto FROM movimiento_ingreso
			GROUP BY agen,fech
			)
			AS ingreso ON
			ingreso.agen=z.obse  AND ingreso.fech=` + fechaAux + `

			-- EGRESO
			LEFT JOIN (
			SELECT agen, fech, SUM(mont) AS monto FROM movimiento_egreso
			GROUP BY agen,fech
			)
			AS egreso ON
			egreso.agen=z.obse  AND egreso.fech=` + fechaAux + `

			-- PRESTAMOS
			LEFT JOIN (
			SELECT agen, fech, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
			GROUP BY agen,fech)
			AS prestamo ON
			prestamo.agen=z.obse AND prestamo.fech=` + fechaAux + `

			-- VIENEN
			LEFT JOIN cobrosypagos ON cobrosypagos.fech=` + fechaAux + `
			AND cobrosypagos.oida=z.oid
			ORDER BY z.obse
			) AS A
	`
	return
}

//generarCobrosYPagosAgencia Generar SQL Cobros y Pagos por Agenacia
func generarCobrosYPagosAgencia(data Pago) (s string) {
	var fecha string
	if data.Desde != "" {
		fecha = ` AND lotepar.fech BETWEEN '` + data.Desde + ` 00:00:00'::TIMESTAMP AND '` + data.Hasta + ` 23:59:59'::TIMESTAMP`
	}
	s = `SELECT cpc.oid, cpc.fech,
				cyp.vien,
				saldo_agencia.saldo,
				debe.monto AS entregado,
				haber.monto AS recibido,
				ingreso.monto AS ingreso,
				egreso.monto AS egreso,
				prestamo.monto AS prestamo,
				prestamo.cuota AS cuota,
				COALESCE(saldo,0) + COALESCE(cyp.vien,0) +
				(
					COALESCE(debe.monto,0) - COALESCE(haber.monto,0)) +
					(
						COALESCE(egreso.monto,0) - (COALESCE(ingreso.monto,0)+COALESCE(prestamo.monto,0))
					) AS van,
				cpc.esta
				FROM cobrosypagoscierre cpc


				--) AS f
				--ON cpc.fech=f.fech

				-- VIENEN
				LEFT JOIN cobrosypagos cyp ON cyp.fech=cpc.fech
				INNER JOIN agencia ON cyp.oida=agencia.oid


				-- DEBE
				LEFT JOIN (
						SELECT agen, fapr AS fapr, SUM(mont) AS monto FROM debe
						GROUP BY agen,fapr
				) AS debe ON
				debe.agen=agencia.obse AND debe.fapr=cyp.fech

				-- HABER
				LEFT JOIN (
					SELECT agen, fapr, SUM(mont) AS monto FROM haber
					GROUP BY agen,fapr
				) AS haber ON
				haber.agen=agencia.obse  AND haber.fapr=cyp.fech

				--INGRESO
				LEFT JOIN (
					SELECT agen, fech,fapr, SUM(mont) AS monto FROM movimiento_ingreso
					GROUP BY agen,fech,fapr
				)
				AS ingreso ON
				ingreso.agen=agencia.obse  AND ingreso.fapr=cyp.fech

				-- EGRESO
				LEFT JOIN (
					SELECT agen, fech,fapr, SUM(mont) AS monto FROM movimiento_egreso
					GROUP BY agen,fech,fapr
				)
				AS egreso ON
				egreso.agen=agencia.obse  AND egreso.fapr=cyp.fech

				-- PRESTAMOS
				LEFT JOIN (
					SELECT agen, fech, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
					GROUP BY agen,fech)
				AS prestamo ON
				prestamo.agen=agencia.obse AND prestamo.fech=cyp.fech


				LEFT JOIN (
					--SELECT saldo_agencia.oid, saldo_agencia.fech, saldo_agencia.saldo,
					--	debe.monto AS entregado, haber.monto AS recibido,
					--	ingreso.monto AS ingreso, egreso.monto AS egreso,
					--	prestamo.monto AS prestamo,prestamo.cuota AS cuota
					--FROM (
					SELECT agencia.oid,agencia.obse, lotepar.fech, SUM(lotepar.saldo) AS saldo
					FROM agencia
					JOIN zr_agencia ON agencia.oid=zr_agencia.oida
					JOIN (
						SELECT arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria
						UNION
						SELECT arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
					) AS lotepar ON zr_agencia.codi=lotepar.agen

					WHERE agencia.obse='` + data.Agencia + `'  ` + fecha + `
					GROUP BY agencia.oid,agencia.obse,lotepar.fech
				) AS saldo_agencia ON saldo_agencia.obse=agencia.obse AND saldo_agencia.fech=cyp.fech


				WHERE agencia.obse='` + data.Agencia + `' AND cyp.oida=agencia.oid
				ORDER BY fech`
	return
}

//GenerarCobrosYPagosSistemas Generar Reporte de pagos por programas Ej: Maticlot, Parley
func (p *Pago) GenerarCobrosYPagosSistemas(data Pago) (jSon []byte, err error) {
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

//generarCobrosYPagosSistemas Creacion de SQL para Reporte de pago  por programas
func generarCobrosYPagosSistemas(data Pago) (s string) {
	var fecha string
	fecha = ` AND lotepar.fech BETWEEN '` + data.Desde + ` 00:00:00'::TIMESTAMP AND '` + data.Hasta + ` 23:59:59'::TIMESTAMP`
	if data.Desde != "" {
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

//GenerarCobrosYPagosDetallados Generar Cobros y Pagos detallados de las ventas
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

//generarCobrosYPagosDetallados Generar SQL para el detalle de las Ventas
func generarCobrosYPagosDetallados(data Pago) (s string) {
	var fecha string
	var sistema string
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

//CobrosYPagos Sentencia de Control
type CobrosYPagos struct {
	Nombre              string
	Venta               float64
	Premio              float64
	Comision            float64
	SaldoAnterio        float64
	Saldo               float64
	Movimiento          float64
	EntregadosRecibidos float64
	Queda               float64
	Participacion       float64
	Total               float64
	Loteria             float64
	Parley              float64
	Calculo             int
	Frecuencia          int
	Sistema             int
	Archivo             int
}

//GenerarCobrosYPagosGrupo Consultando por grupos
func (p *Pago) GenerarCobrosYPagosGrupo() (jSon []byte, err error) {
	//fecha control
	//var fecha string = time.Now().String()[0:10]
	//fmt.Println(p.Cierre)
	s := gCPGrupoDiario(p.Fecha)
	row, err := sys.PostgreSQL.Query(s)
	//fmt.Println(s)
	if err != nil {
		return
	}

	var i, oidAuxiliar, calculo, frecuencia int

	lst := make(map[int]interface{})
	var pago Pago
	var venta, premio, comision, comisionb, saldo, saldovan float64

	for row.Next() {
		var oid, esta int
		var nomb, fvien sql.NullString

		var vent, prem, comi, comical, sald, vien, van float64
		var entr, reci, ingr, egre, pres, cuot, cypsaldo float64

		i++
		e := row.Scan(&oid, &nomb, &calculo, &frecuencia, &vent, &prem,
			&comi, &comical, &sald,
			&entr, &reci, &ingr, &egre, &pres, &cuot,
			&vien, &van, &esta, &fvien, &cypsaldo)

		if e != nil {
			fmt.Println(e.Error())
			return
		}
		if i == 1 {
			oidAuxiliar = oid
		}

		if oidAuxiliar != oid {
			if comision == 0 {
				pago.Comision = comisionb
				pago.Saldo = saldo - pago.Comision
			}
			er := pago.Entregado - pago.Recibido
			movimiento := pago.Egreso - pago.Ingreso
			pago.Van = pago.Saldo + pago.Vienen + er + movimiento
			if p.Cierre == 1 {
				saldo := strconv.FormatFloat(pago.Saldo, 'f', 6, 64)
				monto := strconv.FormatFloat(pago.Van, 'f', 6, 64)
				idgrupo := strconv.Itoa(oidAuxiliar)
				movi := strconv.FormatFloat(movimiento, 'f', 6, 64)
				erec := strconv.FormatFloat(er, 'f', 6, 64)
				s := `UPDATE cobrosypagos_grupo SET sald=` + saldo + `, van=` + monto + `,
				movi=` + movi + `,
				erec=` + erec + `
				WHERE oidg=` + idgrupo + ` AND fech='` + pago.Fecha + `'`
				_, err := sys.PostgreSQL.Exec(s)
				if err != nil {
					fmt.Println(s)
					fmt.Println(err.Error())
				}

				s = `INSERT INTO cobrosypagos_grupo (oidg, fech, vien, sald)
				VALUES (` + idgrupo + `,'` + p.Fecha + ` 00:00:00'::TIMESTAMP + '1 day',` + monto + `,0); `
				_, err = sys.PostgreSQL.Exec(s)
				if err != nil {

					fmt.Println(err.Error())
				}
			}
			lst[oidAuxiliar] = pago
			oidAuxiliar = oid
			venta, premio, comision, comisionb, saldo, saldovan = 0, 0, 0, 0, 0, 0
		}
		venta += vent
		premio += prem
		comisionb += comi
		comision += comical
		saldo += sald
		saldovan += van
		pago.Estatus = esta
		pago.Observacion = util.ValidarNullString(nomb)
		auxiliarFecha := util.ValidarNullString(fvien)
		fmt.Println(auxiliarFecha)
		if auxiliarFecha != "" {

			auxiliarFecha = auxiliarFecha[0:10]
		} else {

			pago.Estatus = 1
		}
		pago.Fecha = auxiliarFecha
		pago.Entregado = entr
		pago.Recibido = reci
		pago.Ingreso = ingr
		pago.Egreso = egre
		pago.Prestamo = pres
		pago.Cuota = cuot
		pago.Vienen = vien

		//pago.Van = saldovan

		pago.Venta = venta
		pago.Premio = premio
		pago.Comision = comision
		pago.Saldo = saldo

	}

	if comision == 0 {
		pago.Comision = comisionb
		pago.Saldo = saldo - pago.Comision
	}
	er := pago.Entregado - pago.Recibido
	movimiento := pago.Egreso - pago.Ingreso
	pago.Van = pago.Saldo + pago.Vienen + er + movimiento

	if p.Cierre == 1 {
		saldo := strconv.FormatFloat(pago.Saldo, 'f', 6, 64)
		monto := strconv.FormatFloat(pago.Van, 'f', 6, 64)
		idgrupo := strconv.Itoa(oidAuxiliar)
		movi := strconv.FormatFloat(movimiento, 'f', 6, 64)
		erec := strconv.FormatFloat(er, 'f', 6, 64)
		s := `UPDATE cobrosypagos_grupo SET sald=` + saldo + `, van=` + monto + `,
		movi=` + movi + `,
		erec=` + erec + `
		WHERE oidg=` + idgrupo + ` AND fech='` + pago.Fecha + `'`
		_, err := sys.PostgreSQL.Exec(s)
		if err != nil {
			//fmt.Println(s)
			fmt.Println(err.Error())
		}

		s = `INSERT INTO cobrosypagos_grupo (oidg, fech, vien, sald)
		VALUES (` + idgrupo + `,'` + p.Fecha + ` 00:00:00'::TIMESTAMP + '1 day',` + monto + `,0); `
		_, err = sys.PostgreSQL.Exec(s)
		if err != nil {

			fmt.Println(err.Error())
		}
	}
	lst[oidAuxiliar] = pago

	jSon, _ = json.Marshal(lst)
	return
}

//gCPGrupoDiario Diario
func gCPGrupoDiario(fecha string) (s string) {
	restar := fecha
	fecha = `'` + fecha + ` 00:00:00'::TIMESTAMP `

	s = `
			-- ##################################
			-- CONSULTA POR GRUPOS SALDOS DIARIOS
			-- ##################################

			SELECT
				f.oid,f.obse, --f.lote,f.parl,f.qued,f.part,
				f.calc,f.freq,
				COALESCE(venta,0) AS venta,
				COALESCE(premio,0) AS premio,
				COALESCE(comision,0) AS comision,
				--venta,premio,comision,
				comisioncal,
				COALESCE((venta-premio-comisioncal),0) AS saldo,
				-- f.soid, --f.slote, f.sparl,f.squed, f.spart,
				-- f.arch,f.fapr,
				entregado,recibido,ingreso,egreso,prestamo, cuota,
				COALESCE(vien,0) AS vien,
				COALESCE(COALESCE((venta-premio-comisioncal),0) + vien + (entregado - recibido) + (ingreso-egreso+prestamo),0) AS van,
				COALESCE(cpc.esta,0) AS esta,cyp.fech, COALESCE(cyp.sald,0) AS cypsaldo
				FROM

				(
				SELECT
					g.oid,g.obse,g.lote,g.parl,g.qued,g.part,g.calc,g.freq, venta,premio,comision,
					CASE
						WHEN g.lote > 0 then (venta * (g.lote + g.parl))/100
						WHEN g.parl > 0 then (venta * (g.lote + g.parl))/100
						WHEN zrg.lote > 0 then (venta * (zrg.lote + zrg.parl))/100
						WHEN zrg.parl > 0 then (venta * (zrg.lote + zrg.parl))/100
					ELSE 0
					END AS comisioncal,
					b.soid,
					zrg.lote AS slote,
					zrg.parl AS sparl,
					zrg.qued AS squed,
					zrg.part AS spart,
					b.arch,
					egreso.fapr,
					COALESCE(debe.monto,0) AS entregado,
					COALESCE(haber.monto,0) AS recibido,
					COALESCE(ingreso.monto,0) AS ingreso,
					COALESCE(egreso.monto,0) AS egreso,
					COALESCE(prestamo.monto,0) AS prestamo,
					COALESCE(prestamo.cuota,0) AS cuota
				FROM
				grupo g LEFT JOIN
				(
				SELECT  g.oid As goid,  SUM(vent) AS venta, SUM(prem) AS premio, SUM(comi) AS comision, s.oid as soid, s.arch FROM (
					SELECT agen, fech, vent,prem,comi, sist from loteria
					UNION
					SELECT agen, fech, vent,prem,comi, sist from parley
				) AS A
				JOIN zr_agencia zr ON A.agen=zr.codi
				JOIN agencia ON agencia.oid = zr.oida
				JOIN grupo g ON g.oid=zr.grupo
				JOIN sistema s ON s.oid=A.sist

				WHERE A.fech = ` + fecha + `
				--AND g.freq=4
				--AND g.obse != 'AGE. DIRECTAS'
				GROUP BY g.oid, s.oid,  s.arch
				ORDER BY g.oid
				) AS b ON g.oid=B.goid
				LEFT JOIN zr_negociacion_grupo zrg ON zrg.oids=b.soid AND g.oid=zrg.oidg

				-- DEBE
				LEFT JOIN (
					SELECT grupo, fapr, SUM(mont) AS monto FROM debe
					GROUP BY grupo,fapr
				) AS debe ON
				debe.grupo=g.oid AND debe.fapr=` + fecha + `

				-- HABER
				LEFT JOIN (
				SELECT grupo, fapr, SUM(mont) AS monto FROM haber
				GROUP BY grupo,fapr
				) AS haber ON
				haber.grupo=g.oid AND haber.fapr=` + fecha + `

				--INGRESO
				LEFT JOIN (
				SELECT grupo, fech,fapr, SUM(mont) AS monto FROM movimiento_ingreso
				GROUP BY grupo,fech,fapr
				)
				AS ingreso ON
				ingreso.grupo=g.oid AND ingreso.fapr=` + fecha + `

				-- EGRESO
				LEFT JOIN (
				SELECT grupo, fech,fapr, SUM(mont) AS monto FROM movimiento_egreso
				GROUP BY grupo,fech,fapr
				)
				AS egreso ON
				egreso.grupo=g.oid AND egreso.fapr=` + fecha + `

				-- PRESTAMOS
				LEFT JOIN (
					SELECT grupo, fapr, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
					GROUP BY grupo,fapr)
				AS prestamo ON
				prestamo.grupo=g.oid AND prestamo.fapr=` + fecha + `
				) AS f
			LEFT JOIN cobrosypagos_grupo cyp ON cyp.fech='` + restar + ` 00:00:00'::TIMESTAMP  AND cyp.oidg=f.oid
			LEFT JOIN cobrosypagoscierre_grupo cpc ON cpc.fech='` + restar + ` 00:00:00'::TIMESTAMP
			--WHERE f.obse != '0'
			ORDER BY f.oid
			-- + '-24:00:00'
			--
	`

	return
}

func gCPGrupoParticipacionDiaria(fecha string) (s string) {
	sumar := `'` + fecha + ` 00:00:00'::TIMESTAMP + '1 day'`
	s = `
				-- ##################################
				-- CONSULTA POR GRUPOS SALDOS DIARIOS
				-- ##################################
				INSERT INTO movimiento_egreso (comer,grupo,subgr,colec,agenc,fech,fapr,freg,tipo,cuen, oper,obse,mont)
				SELECT * FROM (
				SELECT 1, P.oid, 0,0,0,'` + fecha + `'::DATE,` + sumar + `,now(), 1, 0, '', 'PAGO POR PARTICIPACION - ` + fecha + `', COALESCE(SUM(participacion),0) AS part FROM (
				SELECT
				f.oid,f.obse,agencianombre,--f.lote,f.parl,f.qued,f.part,
				f.calc,f.freq,
				venta,premio,comision, comisioncal, (venta-premio-comisioncal) AS saldo,
				f.soid, --f.slote, f.sparl,
				f.squed, f.spart, ((venta-premio-comisioncal)* f.spart)/100 AS participacion,
				-- f.arch,f.fapr,
				entregado,recibido,ingreso,egreso,prestamo, cuota,
				vien,
				(venta-premio-comisioncal) + vien + (entregado - recibido) + (ingreso-egreso+prestamo) AS van,
				COALESCE(cpc.esta,0) AS esta,cyp.fech
				FROM

				(
				SELECT
					g.oid,g.obse,agencianombre,g.lote,g.parl,g.qued,g.part,g.calc,g.freq,
					COALESCE(venta,0) AS venta,
					COALESCE(premio,0) AS premio,
					COALESCE(comision,0) AS comision,
					CASE
						WHEN g.lote > 0 then (venta * (g.lote + g.parl))/100
						WHEN g.parl > 0 then (venta * (g.lote + g.parl))/100
						WHEN zrg.lote > 0 then (venta * (zrg.lote + zrg.parl))/100
						WHEN zrg.parl > 0 then (venta * (zrg.lote + zrg.parl))/100
					ELSE 0
					END AS comisioncal,
					b.soid,
					zrg.lote AS slote,
					zrg.parl AS sparl,
					zrg.qued AS squed,
					zrg.part AS spart,
					b.arch,
					egreso.fapr,
					COALESCE(debe.monto,0) AS entregado,
					COALESCE(haber.monto,0) AS recibido,
					COALESCE(ingreso.monto,0) AS ingreso,
					COALESCE(egreso.monto,0) AS egreso,
					COALESCE(prestamo.monto,0) AS prestamo,
					COALESCE(prestamo.cuota,0) AS cuota

				FROM
				grupo g LEFT JOIN
				(
					SELECT  g.oid As goid, agencia.obse AS agencianombre,  SUM(vent) AS venta, SUM(prem) AS premio, SUM(comi) AS comision, s.oid as soid, s.arch FROM (
						SELECT agen, fech, vent,prem,comi, sist from loteria
						UNION
						SELECT agen, fech, vent,prem,comi, sist from parley
					) AS A
					JOIN zr_agencia zr ON A.agen=zr.codi
					JOIN agencia ON agencia.oid = zr.oida
					JOIN grupo g ON g.oid=zr.grupo
					JOIN sistema s ON s.oid=A.sist

					WHERE A.fech ='` + fecha + `'
					--AND g.freq=1
					--AND g.obse != 'AGE. DIRECTAS'
					GROUP BY g.oid, agencia.obse, s.oid,  s.arch
					ORDER BY g.oid
				) AS b ON g.oid=B.goid
				LEFT JOIN zr_negociacion_grupo zrg ON zrg.oids=b.soid AND g.oid=zrg.oidg

				-- DEBE
				LEFT JOIN (
						SELECT grupo, fapr, SUM(mont) AS monto FROM debe
						GROUP BY grupo,fapr
				) AS debe ON
				debe.grupo=g.oid AND debe.fapr='` + fecha + ` 00:00:00'::TIMESTAMP

				-- HABER
				LEFT JOIN (
					SELECT grupo, fapr, SUM(mont) AS monto FROM haber
					GROUP BY grupo,fapr
				) AS haber ON
				haber.grupo=g.oid AND haber.fapr='` + fecha + ` 00:00:00'::TIMESTAMP

				--INGRESO
				LEFT JOIN (
					SELECT grupo, fech,fapr, SUM(mont) AS monto FROM movimiento_ingreso
					GROUP BY grupo,fech,fapr
				)
				AS ingreso ON
				ingreso.grupo=g.oid AND ingreso.fapr='` + fecha + ` 00:00:00'::TIMESTAMP

				-- EGRESO
				LEFT JOIN (
					SELECT grupo, fech,fapr, SUM(mont) AS monto FROM movimiento_egreso
					GROUP BY grupo,fech,fapr
				)
				AS egreso ON
				egreso.grupo=g.oid AND egreso.fapr='` + fecha + ` 00:00:00'::TIMESTAMP

				-- PRESTAMOS
				LEFT JOIN (
					SELECT grupo, fapr, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
					GROUP BY grupo,fapr)
				AS prestamo ON
				prestamo.grupo=g.oid AND prestamo.fapr='` + fecha + ` 00:00:00'::TIMESTAMP


				) AS f
				LEFT JOIN cobrosypagos_grupo cyp ON cyp.fech='` + fecha + ` 00:00:00'::TIMESTAMP  AND cyp.oidg=f.oid
				LEFT JOIN cobrosypagoscierre_grupo cpc ON cpc.fech='` + fecha + ` 00:00:00'::TIMESTAMP
				WHERE f.freq=1
				) AS P -- PARTICIPACION
				GROUP BY P.oid ) AS B
				WHERE B.part>0
	`
	return
}

//GenerarCobrosYPagosGMQ Consultando por grupos
func (p *Pago) GenerarCobrosYPagosGMQ(data Pago) (jSon []byte, err error) {
	//fecha control
	var fecha string = time.Now().String()[0:10]

	row, err := sys.PostgreSQL.Query(gCPGrupoMensual(fecha))

	if err != nil {
		return
	}
	var oidAuxiliar int
	lst := make(map[string]interface{})
	CalculoData := gCPGrupoMensualQueda("")
	for row.Next() {
		var oid, calculo, fecuencia, i int

		var cobros CobrosYPagos
		var aoid string
		var loteria, parley, queda, parti, venta, premio, comision float64
		i++
		e := row.Scan(&oid, &aoid, &loteria, &parley,
			&queda, &parti, &calculo, &fecuencia, &venta, &premio, &comision)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		if i == 1 {
			oidAuxiliar = oid
		}

		if oidAuxiliar != oid {

		}
		cobros.Nombre = aoid
		cobros.Calculo = calculo
		cobros.Frecuencia = fecuencia
		cobros.Loteria, cobros.Parley, cobros.Queda = obtenerSaldoGMQ(CalculoData, strconv.Itoa(oid))
		cobros.Saldo = cobros.Loteria + cobros.Parley
		lst[strconv.Itoa(oid)] = cobros

	}

	//fmt.Println(CalculoQueda["496"])
	jSon, _ = json.Marshal(lst)
	return
}

//obtenerSaldoGMQ obtiene el saldo total del grupo aplicando las reglas de negocio
func obtenerSaldoGMQ(lst map[string][]CobrosYPagos, valorDeseao string) (saldol float64, saldop float64, queda float64) {
	for c, v := range lst {
		if c == valorDeseao {
			for _, vl := range v {
				var comision float64
				var calc float64
				if vl.Archivo == 0 {
					comision = (vl.Venta * vl.Loteria) / 100
					calc = vl.Venta - vl.Premio - comision
					saldol += calc
				} else {
					comision = (vl.Venta * vl.Parley) / 100
					calc = vl.Venta - vl.Premio - comision
					saldop += calc
				}

				if vl.Queda > 0 {

					if calc > 0 {
						queda += (calc * vl.Queda) / 100
					}
				}

			}

		}
	}
	return
}

//gCPGrupoMensual Diario
func gCPGrupoMensual(fecha string) (s string) {
	/**
	goid,g.obse,g.trip,g.term,g.qued,g.part,g.calc,g.freq,aoid,
	a.obse,a.trip,a.term,a.qued,a.part,a.freq, venta,premio,comision, arch
	**/
	s = `
		-- GLOBAL
		SELECT
		goid,g.obse,g.lote,g.parl,g.qued,g.part,g.calc,g.freq, venta,premio,comision FROM (
			SELECT  g.oid As goid,  SUM(vent) AS venta, SUM(prem) AS premio, SUM(comi) AS comision FROM (
				SELECT agen, fech, vent,prem,comi, sist from loteria
				UNION
				SELECT agen, fech, vent,prem,comi, sist from parley
			) AS A
			JOIN zr_agencia zr ON A.agen=zr.codi
			JOIN agencia ON agencia.oid = zr.oida
			JOIN grupo g ON g.oid=zr.grupo
			JOIN sistema s ON s.oid=A.sist

			WHERE A.fech BETWEEN '2017-01-01 00:00:00'::TIMESTAMP AND '2017-01-31 23:59:59'::TIMESTAMP
			AND g.obse != 'AGE. DIRECTAS'
			GROUP BY g.oid
			ORDER BY g.oid
		) AS b
		JOIN grupo g ON g.oid=B.goid

	`

	return
}

//gCPGrupoMensualQueda Queda
func gCPGrupoMensualQueda(fecha string) (lst map[string][]CobrosYPagos) {
	/**
	,zrg.parl,zrg.qued,zrg.part,zrg.calc,zrg.freq,
	venta,premio,comision, b.sist,b.arch
	**/
	s := `
		-- GLOBAL POR PROGRAMAS
		SELECT goid,g.obse,
			zrg.lote,zrg.parl,zrg.qued,zrg.part,venta,premio,comision, b.sist,b.arch FROM (
			SELECT  g.oid As goid,  A.sist, s.arch, SUM(vent) AS venta, SUM(prem) AS premio, SUM(comi) AS comision FROM (
				SELECT agen, fech, vent,prem,comi, sist from loteria
				UNION
				SELECT agen, fech, vent,prem,comi, sist from parley
			) AS A
			JOIN zr_agencia zr ON A.agen=zr.codi
			JOIN agencia ON agencia.oid = zr.oida
			JOIN grupo g ON g.oid=zr.grupo
			JOIN sistema s ON s.oid=A.sist

			WHERE A.fech BETWEEN '2017-01-01 00:00:00'::TIMESTAMP AND '2017-01-31 23:59:59'::TIMESTAMP
			AND g.freq=4
			AND g.obse != 'AGE. DIRECTAS'


			GROUP BY g.oid, A.sist, s.arch
			ORDER BY g.oid
		) AS b
		JOIN grupo g ON g.oid=b.goid
		LEFT JOIN zr_negociacion_grupo zrg ON zrg.oids=b.sist AND g.oid=zrg.oidg
	`
	row, err := sys.PostgreSQL.Query(s)

	if err != nil {
		return
	}
	var i int
	var oidAuxiliar int
	lst = make(map[string][]CobrosYPagos)
	var cobros []CobrosYPagos
	for row.Next() {
		var oid, sist, arch int
		var aoid string
		var loteria, parley, queda, parti, venta, premi, comi float64
		var cobro CobrosYPagos
		e := row.Scan(&oid, &aoid, &loteria, &parley, &queda, &parti,
			&venta, &premi, &comi, &sist, &arch)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		i++

		if i == 1 {
			oidAuxiliar = oid
		}

		if oidAuxiliar != oid {
			//fmt.Println(oidAuxiliar)
			lst[strconv.Itoa(oidAuxiliar)] = cobros

			oidAuxiliar = oid
			cobros = nil
		}

		cobro.Loteria = loteria
		cobro.Parley = parley
		cobro.Queda = queda
		cobro.Participacion = parti
		cobro.Venta = venta
		cobro.Premio = premi
		cobro.Comision = comi
		cobro.Sistema = sist
		cobro.Archivo = arch
		cobros = append(cobros, cobro)

	}
	lst[strconv.Itoa(oidAuxiliar)] = cobros
	return
}

//EstadoDeCuenta Estructura
type EstadoDeCuenta struct {
	Viene      float64
	Saldo      float64
	Movimiento float64
	Entregado  float64
	Total      float64
}

//EstadoDeCuentasGrupo Control General
func (p *Pago) EstadoDeCuentasGrupo() (jSon []byte, err error) {
	fecha := ` AND cg.fech BETWEEN '` + p.Desde + ` 00:00:00'::TIMESTAMP AND '` + p.Hasta + ` 23:59:59'::TIMESTAMP`
	s := `SELECT fech,vien, sald, movi, erec, van FROM grupo g
				JOIN cobrosypagos_grupo cg ON g.oid=cg.oidg
				where g.oid=` + strconv.Itoa(p.Oid) + fecha + `
				AND van!=0
				ORDER BY cg.fech`
	//fmt.Println(s)

	lst := make(map[string]interface{})
	row, err := sys.PostgreSQL.Query(s)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for row.Next() {
		var fecha string
		var vien, sald, movi, erec, total float64
		var ec EstadoDeCuenta
		e := row.Scan(&fecha, &vien, &sald, &movi, &erec, &total)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		ec.Viene = vien
		ec.Saldo = sald
		ec.Movimiento = movi
		ec.Entregado = erec
		ec.Total = total
		lst[fecha[0:10]] = ec
	}

	jSon, _ = json.Marshal(lst)
	return
}
