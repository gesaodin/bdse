package balance

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/util"
)

type Pago struct {
	Agencia     string  `json:"agencia,omitempty"`
	Voucher     string  `json:"voucher,omitempty"`
	Desde       string  `json:"desde,omitempty"`
	Hasta       string  `json:"hasta,omitempty"`
	Deposito    string  `json:"deposito,omitempty"`
	Banco       int     `json:"banco,omitempty"`
	FormaDePago int     `json:"forma,omitempty"`
	Monto       float64 `json:"monto,omitempty"`
	Saldo       float64 `json:"saldo,omitempty"`
	Vienen      float64 `json:"vienen,omitempty"`
	Recibido    float64 `json:"recibido,omitempty"`
	Entregado   float64 `json:"entregado,omitempty"`
	Ingreso     float64 `json:"ingreso,omitempty"`
	Egreso      float64 `json:"egreso,omitempty"`
	Prestamo    float64 `json:"prestamo,omitempty"`
	Cuota       float64 `json:"cuota,omitempty"`
}

type Respuesta struct {
	Cantidad int64  `json:"cant"`
	Msj      string `json:"msj"`
}

func (p *Pago) Registrar(data Pago) (jSon []byte, err error) {
	monto := strconv.FormatFloat(data.Monto, 'f', 6, 64)
	forma := strconv.Itoa(data.FormaDePago)
	banco := strconv.Itoa(data.Banco)
	// deposito := data.Deposito.String()
	// desde := data.Desde.String()
	// hasta := data.Hasta.String()
	s := "INSERT INTO haber (agen,mont,vouc,fech,fhas,fdep,freg,tipo,banc) VALUES "
	s += "('" + data.Agencia + "'," + monto + ",'" + data.Voucher + "',"
	s += "'" + data.Desde + "','" + data.Hasta + "','" + data.Deposito + "',now(),"
	s += forma + "," + banco + ");"
	// s := "INSERT INTO haber (agen,mont,vouc,fech,fhas,fdep,freg,tipo,banc) VALUES "
	// s += "(''" + data.Agencia + "',''" + monto + "',''" + data.Voucher + "',"
	// s += "'" + desde[0:10] + "','" + hasta[0:10] + "','" + deposito[0:10] + "',now(),"
	// s += forma + "," + banco + ");"
	fmt.Println(s)
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

func (p *Pago) GenerarCobrosYPagos(data Pago) (jSon []byte, err error) {

	s := `
				SELECT t.agen, SUM(t.vent-t.prem-t.comi) AS saldo,
					debe.monto AS entregado, haber.monto AS recibido,
					ingreso.monto AS ingreso, egreso.monto AS egreso,
					prestamo.monto AS prestamo,
					vienen.vien AS vienen, SUM(prestamo.cuota) AS cuota
				FROM
					(
						select arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria UNION
						select arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
					) AS t

				LEFT JOIN (
						SELECT agen, fdep, SUM(mont) AS monto FROM debe
						WHERE fdep BETWEEN '2017-02-01 00:00:00'::TIMESTAMP AND '2017-02-01 23:59:59'::TIMESTAMP GROUP BY agen,fdep
				) AS debe ON
				debe.agen=t.agen

				LEFT JOIN (
					SELECT agen, fdep, SUM(mont) AS monto FROM haber
					WHERE fdep BETWEEN '2017-02-01 00:00:00'::TIMESTAMP AND '2017-02-01 23:59:59'::TIMESTAMP GROUP BY agen,fdep
				) AS haber ON
				haber.agen=t.agen

				LEFT JOIN (
					SELECT agen, fech, SUM(mont) AS monto FROM movimiento_ingreso
					WHERE fech BETWEEN '2017-02-01 00:00:00'::TIMESTAMP AND '2017-02-01 23:59:59'::TIMESTAMP GROUP BY agen,fech
				)
				AS ingreso ON
				ingreso.agen=t.agen

				LEFT JOIN (
					SELECT agen, fech, SUM(mont) AS monto FROM movimiento_egreso
					WHERE fech BETWEEN '2017-02-01 00:00:00'::TIMESTAMP AND '2017-02-01 23:59:59'::TIMESTAMP GROUP BY agen,fech
				)
				AS egreso ON
				egreso.agen=t.agen

				LEFT JOIN (
					SELECT agen, fech, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
					WHERE fech BETWEEN '2017-02-01 00:00:00'::TIMESTAMP AND '2017-02-01 23:59:59'::TIMESTAMP GROUP BY agen,fech
				)
				AS prestamo ON
				prestamo.agen=t.agen

				LEFT JOIN (
						SELECT * FROM cobrosypagos ORDER BY fech ASC LIMIT 1
				) AS vienen
				ON vienen.agen=t.agen

				WHERE t.fech='2017-01-02'
				GROUP by t.agen, debe.monto, haber.monto, ingreso.monto, egreso.monto,  prestamo.monto, vienen.vien
				ORDER by t.agen
			`

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
