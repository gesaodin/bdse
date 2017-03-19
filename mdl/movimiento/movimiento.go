//ingreso y egreso ayudan a las ganancias o perdidas
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
	Monto             float64 `json:"monto,omitempty"`
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

//MSJ Control de mensaje
type MSJ struct {
	Msj  string `json:"msj"`
	Tipo int    `json:"tipo"`
}

//Salvar datos
func (m *Movimiento) Salvar() (jSon []byte, err error) {
	var ingreso, egreso string
	switch m.TipoDeOperacion {
	case 1:
		ingreso, _ = m.generarSQL()
	case 2:
		_, egreso = m.generarSQL()
	default:
		ingreso, egreso = m.generarSQL()
	}

	_, err = sys.PostgreSQL.Exec(ingreso + egreso)

	if err != nil {

		return
	}
	var res MSJ
	res.Msj = "Se inserto correctamente"
	res.Tipo = 1
	jSon, err = json.Marshal(res)
	return
}

//Actualizar datos
func (m *Movimiento) Actualizar() (jSon []byte, err error) {
	tabla := "haber"
	if m.FormaDePago == 0 {
		tabla = "debe"
	}

	s := `UPDATE  ` + tabla + ` SET fapr = '` + m.Fecha + `', resp='` + m.Observacion + `',
	esta=` + strconv.Itoa(m.Estatus) + ` WHERE oid =` + strconv.Itoa(m.Oid)

	_, err = sys.PostgreSQL.Exec(s)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var res MSJ
	res.Msj = "Se actualizo correctamente..."
	res.Tipo = 1
	jSon, err = json.Marshal(res)
	return
}

//generarSQL Consultar
func (m *Movimiento) generarSQL() (sqlI string, sqlE string) {
	sql1 := "INSERT INTO "
	ie := "(comer,grupo,subgr,colec,agenc,fech,freg,tipo,cuen,mont,oper,obse, toke)" // INGRESO | EGRESO

	monto := strconv.FormatFloat(m.Monto, 'f', 2, 64)
	iii := "(" + strconv.Itoa(m.Comercializadora) + "," + strconv.Itoa(m.Grupo)
	iii += "," + strconv.Itoa(m.SubGrupo) + "," + strconv.Itoa(m.Colector) + "," + strconv.Itoa(m.AgenciaCod)
	iii += ",'" + m.Fecha + "',now(),"
	cuenta := strconv.Itoa(m.TipoDebe) + "," + strconv.Itoa(m.CuentaDebe) + ","
	iff := monto + ", '" + m.Voucher + "', '" + m.Observacion + "', md5('" + m.Fecha + m.Voucher + monto + "'));"
	sqlI = sql1 + "movimiento_ingreso " + ie + " VALUES " + iii + cuenta + iff

	cuenta = strconv.Itoa(m.TipoHaber) + "," + strconv.Itoa(m.CuentaHaber) + ","
	sqlE = sql1 + "movimiento_egreso " + ie + " VALUES " + iii + cuenta + iff
	//sqlE = sqls + sqle

	return
}

//Listar todos los movimientos por fechas
func (m *Movimiento) Listar() (jSon []byte, err error) {
	var lst []interface{}
	s := `SELECT A.fech,A.oper, CONCAT(C.nomb, ' ', C.num ) AS debe, A.tipo AS tdebe,
		CONCAT(D.nomb, ' ', D.num ) AS haber, B.tipo AS thaber, A.obse, A.mont, A.toke
		 FROM movimiento_ingreso A
		INNER JOIN movimiento_egreso B ON A.toke=B.toke
		JOIN cuenta  C ON C.cod = A.cuen
		JOIN cuenta D ON D.cod = B.cuen  WHERE A.fech='` + m.Fecha + `';` //` +  AND fapr != ''; m.FDeposito + `

	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for row.Next() {
		var movimiento Movimiento
		var fech, toke, debe, haber string

		var obse, oper sql.NullString

		var tdebe, thaber int
		var mont sql.NullFloat64

		e := row.Scan(&fech, &oper, &debe, &tdebe, &haber, &thaber, &obse, &mont, &toke)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		movimiento.Fecha = fech
		movimiento.Voucher = util.ValidarNullString(oper)
		movimiento.Observacion = util.ValidarNullString(obse)
		movimiento.Monto = util.ValidarNullFloat64(mont)
		movimiento.CuentaDebeNombre = debe
		movimiento.TipoDebe = tdebe
		movimiento.CuentaHaberNombre = haber
		movimiento.TipoHaber = thaber
		movimiento.Token = toke
		lst = append(lst, movimiento)
	}
	jSon, err = json.Marshal(lst)
	return
}

//ListarDepositos Depositos Pendientes
func (m *Movimiento) ListarDepositos() (jSon []byte, err error) {
	var lst []interface{}
	s := `SELECT banco.nomb, debe.oid,agen,debe.mont,vouc,fdep,tipo,banc,resp FROM debe
	LEFT JOIN banco ON debe.banc=banco.oid	WHERE esta=0` //` +  AND fapr != ''; m.FDeposito + `
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for row.Next() {
		var movimiento Movimiento
		var agen, fdep, vouc string
		var nomb, resp sql.NullString

		var oid, tipo, banc int
		var mont sql.NullFloat64

		e := row.Scan(&nomb, &oid, &agen, &mont, &vouc, &fdep, &tipo, &banc, &resp)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		movimiento.Oid = oid
		movimiento.Agencia = agen
		//movimiento.FDeposito = fdep
		movimiento.Voucher = vouc
		movimiento.Observacion = util.ValidarNullString(resp)
		movimiento.Banco = banc
		movimiento.BancoNombre = util.ValidarNullString(nomb)
		movimiento.Monto = util.ValidarNullFloat64(mont)
		lst = append(lst, movimiento)
	}
	jSon, err = json.Marshal(lst)
	return
}

//ListarCuentas 0: Cuentas y 1: Banco
func (m *Movimiento) ListarCuentas() (jSon []byte, err error) {
	var lst []interface{}

	s := `SELECT cod,nomb,num, tipo FROM cuenta `
	if m.TipoDeOperacion == 1 {
		s = `SELECT cod, nomb, num, tipo FROM cuenta WHERE esta=1`
	}

	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for row.Next() {
		var movimiento Movimiento
		var cuen, num sql.NullString
		var cod, tipo int
		e := row.Scan(&cod, &cuen, &num, &tipo)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		movimiento.Oid = cod
		movimiento.Nombre = util.ValidarNullString(cuen) + " " + util.ValidarNullString(num)
		movimiento.TipoDeOperacion = tipo
		lst = append(lst, movimiento)
	}
	jSon, err = json.Marshal(lst)
	return
}

//
