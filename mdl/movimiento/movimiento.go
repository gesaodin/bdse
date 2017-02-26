package movimiento

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/util"
)

type Movimiento struct {
	Oid             int     `json:"oid,omitempty"`
	Agencia         string  `json:"agencia,omitempty"`
	Fecha           string  `json:"fecha,omitempty"`
	FDeposito       string  `json:"fdeposito,omitempty"`
	FOperacion      string  `json:"foperacion,omitempty"`
	Voucher         string  `json:"voucher,omitempty"`
	FormaDePago     int     `json:"forma,omitempty"`
	TipoDeOperacion int     `json:"operacion,omitempty"`
	TipoTabla       int     `json:"tipo,omitempty"`
	Monto           float64 `json:"monto,omitempty"`
	Cuota           float64 `json:"cuota,omitempty"`
	Cuenta          int     `json:"cuenta,omitempty"`
	Banco           int     `json:"banco,omitempty"`
}

type MSJ struct {
	Msj  string `json:"msj"`
	Tipo int    `json:"tipo"`
}

func (m *Movimiento) Salvar() (jSon []byte, err error) {

	_, err = sys.PostgreSQL.Exec(m.generarSQL())

	if err != nil {

		return
	}
	var res MSJ
	res.Msj = "Se inserto correctamente"
	res.Tipo = 1
	jSon, err = json.Marshal(res)
	return
}

//
func (m *Movimiento) generarSQL() (sql string) {
	sql = "INSERT INTO "
	ie := "(agen,fech,tipo,cuen,banc,	form,	mont)"               // INGRESO | EGRESO
	dh := "(agen,mont,vouc,fdep,freg,fope,fapro,tipo,banc,esta)" //DEBE | HABER
	monto := strconv.FormatFloat(m.Monto, 'f', 2, 64)
	iie := "('" + m.Agencia + "','" + m.Fecha + "'," + strconv.Itoa(m.Cuenta) + ","
	iie += strconv.Itoa(m.Banco) + "," + strconv.Itoa(m.FormaDePago) + ","
	iie += monto + ")"
	operacion := m.FOperacion
	if m.FOperacion == "" {
		operacion = m.FDeposito
	}
	idh := "('" + m.Agencia + "'," + monto + ",'" + m.Voucher
	idh += "','" + m.FDeposito + "',Now(),'" + operacion + "','" + operacion + "',"
	idh += strconv.Itoa(m.FormaDePago) + "," + strconv.Itoa(m.Banco)
	idh += "," + strconv.Itoa(m.TipoDeOperacion) + ")"

	switch m.TipoTabla {
	case 0:

		sql += "movimiento_egreso " + ie + " VALUES " + iie
		break
	case 1:
		sql += "movimiento_ingreso " + ie + " VALUES " + iie
		break
	case 2:
		cuota := strconv.FormatFloat(m.Cuota, 'f', 2, 64)
		t := "(agen,tipo,fech,mcuo,cuen,saldo,banc,form,mont)"
		p := "('" + m.Agencia + "'," + strconv.Itoa(m.TipoDeOperacion) + ",'"
		p += m.Fecha + "'," + cuota + "," + strconv.Itoa(m.Cuenta) + ","
		p += monto + "," + strconv.Itoa(m.Banco) + "," + strconv.Itoa(m.FormaDePago)
		p += "," + monto

		sql += "movimiento_prestamo " + t + " VALUES " + p
		break
	case 3:
		sql += "debe " + dh + " VALUES " + idh
		break
	case 4:
		sql += "haber " + dh + " VALUES " + idh
		break
	default:
		sql = ""
	}
	return

}

func (m *Movimiento) ListarDepositos() (jSon []byte, err error) {
	var lst []interface{}
	s := `SELECT oid,agen,mont,vouc,fdep,tipo,banc FROM haber WHERE esta=0
  AND fdep='` + m.FDeposito + `';`
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for row.Next() {
		var movimiento Movimiento
		var agen, fdep, vouc string
		var oid, tipo, banc int
		var mont sql.NullFloat64

		row.Scan(&oid, &agen, &mont, &vouc, &fdep, &tipo, &banc)
		movimiento.Oid = oid
		movimiento.Agencia = agen
		//movimiento.FDeposito = fdep
		movimiento.Voucher = vouc
		movimiento.Banco = banc
		movimiento.Monto = util.ValidarNullFloat64(mont)
		lst = append(lst, movimiento)
	}
	jSon, err = json.Marshal(lst)
	return
}

//
