package util

import "strconv"

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
