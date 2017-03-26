//El hombre que busca la verdad para compartirla se hace mejor fc

//hasta los 39 sere petulante y sobervio y orgulloso porque a partir de los 40 seré perfecto
//Un negro en la nieve es un blanco perfecto
package transferencia

import (
  "database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gesaodin/bdse/sys"
  "github.com/gesaodin/bdse/util"
)


//Pago Control de Pagos
type Transferencia struct {
	Oid             int     `json:"oid,omitempty"`
	Banca           int     `json:"banca,omitempty"`
	Grupo           int     `json:"grupo,omitempty"`
	SubGrupo        int     `json:"subgrupo,omitempty"`
	Colector        int     `json:"colector,omitempty"`
	Agencia         string  `json:"agencia,omitempty"`
  Nombre          string  `json:"nombre,omitempty"`
  Cedula          string  `json:"cedula,omitempty"`
  RazonSocial     string  `json:"razon,omitempty"`
  CuentaBancaria  string  `json:"cuenta,omitempty"`
  Correo          string  `json:"correo,omitempty"`
  Ticket          string  `json:"ticket,omitempty"`
  Serial          string  `json:"serial,omitempty"`
  Sistema         int     `json:"sistema,omitempty"`
  NombreSistema   string  `json:"nombresis,omitempty"`
  MontoTicket     float64 `json:"montot,omitempty"`
  MontoSolicitado float64 `json:"montos,omitempty"`
  Fecha           string  `json:"fecha,omitempty"`
  Estatus         int     `json:"estatus,omitempty"`
}

//Respuesta Generales
type Respuesta struct {
	Cantidad int64  `json:"cant"` // Cantidad de elementos
	Msj      string `json:"msj"`  // Mensaje almacenado
}

//Registrar una solicitud de transferencia
func (t *Transferencia) Registrar() (jSon []byte, err error){
  montot := strconv.FormatFloat(t.MontoTicket, 'f', 6, 64)
  montos := strconv.FormatFloat(t.MontoSolicitado, 'f', 6, 64)

	s := "INSERT INTO solicitud_transferencia (comer,grupo,subgr,colec,oida,cedul,nombr,corre,cuent,ticke,seria,sist,montt,monts,fech,esta) VALUES "
	s += "(" + strconv.Itoa(t.Banca) + "," + strconv.Itoa(t.Grupo) + "," + strconv.Itoa(t.SubGrupo)
	s += "," + strconv.Itoa(t.Colector) + "," +  strconv.Itoa(t.Oid)
	s += ",'" + t.Cedula + "','" + t.RazonSocial + "','" + t.CuentaBancaria + "',"
	s += "'" + t.Correo + "','" + t.Ticket + "','" + t.Serial
	s += "'," + strconv.Itoa(t.Sistema) + "," + montot + "," + montos + ",now(), 0);"

  //	fmt.Println(s)
	rs, err := sys.PostgreSQL.Exec(s)
	if err != nil {
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


func (t *Transferencia) ListarAgencia() (jSon []byte, err error){
  var oida string
  if t.Oid != 0 {
    oida = "AND st.oida=" + strconv.Itoa(t.Sistema)
  }
  s := `SELECT
          oida,agencia.obse, sistema.obse AS sistema, cedul,nombr,corre,cuent,ticke,seria,sist,montt,monts,fech, esta
        FROM solicitud_transferencia st
        JOIN agencia ON agencia.oid=st.oida
        JOIN sistema ON sistema.oid=st.sist
        WHERE st.comer=0 AND st.grupo=0 AND st.subgr=0 AND st.colec=0
         AND st.esta=0 ` + oida

  fmt.Println(s)
  row, err := sys.PostgreSQL.Query(s)
	if err != nil {
    fmt.Println(s)
		return
	}


  var lst []interface{}
	for row.Next() {
    var t Transferencia
		var oid, sist, esta int
    var cedu, agencia,sistema, nomb, corr, cuen, tick, seri, fech sql.NullString
    var fecha string
		var montt, montos sql.NullFloat64


		e := row.Scan(&oid, &agencia, &sistema, &cedu, &nomb, &corr, &cuen, &tick, &seri, &sist, &montt, &montos, &fech, &esta)
		if e != nil {
			fmt.Println(e.Error())
			return
		}

    if util.ValidarNullString(fech) != "" {
      fecha = util.ValidarNullString(fech)[0:10]
    }

    t.Oid = oid
    t.Nombre = util.ValidarNullString(agencia)
		t.Cedula = util.ValidarNullString(cedu)
		t.RazonSocial = util.ValidarNullString(nomb)
    t.CuentaBancaria = util.ValidarNullString(cuen)
    t.Correo = util.ValidarNullString(corr)
    t.Ticket = util.ValidarNullString(tick)
    t.Serial = util.ValidarNullString(seri)
    t.Sistema = sist
    t.NombreSistema = util.ValidarNullString(sistema)
    t.MontoTicket = util.ValidarNullFloat64(montt)
    t.MontoSolicitado = util.ValidarNullFloat64(montos)
    t.Fecha = fecha

		lst = append(lst, t)
	}

	jSon, _ = json.Marshal(lst)
  return
}


func (t *Transferencia) ListarGrupo() (jSon []byte, err error){
  s := `SELECT
          oidg,grupo.obse, sistema.obse AS sistema,cedul,nombr,corre,cuent,ticke,seria,sist,montt,monts,fech, esta
        FROM solicitud_transferencia st
        JOIN grupo ON grupo.oid=st.oida
        JOIN sistema ON sistema.oid=st.sist
        WHERE
          st.comer=0 AND st.subgr=0 AND st.colec=0 AND oida=0 AND st.esta=0`

  row, err := sys.PostgreSQL.Query(s)
  if err != nil {
    return
  }


  var lst []interface{}
  for row.Next() {
    var t Transferencia
    var oid, sist, esta int
    var grupo, sistema, cedu, nomb, corr, cuen, tick, seri, fech sql.NullString
    var fecha string
    var montt, montos sql.NullFloat64


    e := row.Scan(&oid, &grupo, &sistema, &cedu, &nomb, &corr, &cuen, &tick, &seri, &sist, &montt, &montos, &fech, &esta)
    if e != nil {
      fmt.Println(e.Error())
      return
    }

    if util.ValidarNullString(fech) != "" {
      fecha = util.ValidarNullString(fech)[0:10]
    }

    t.Oid = oid
    t.Nombre = util.ValidarNullString(grupo)
    t.Cedula = util.ValidarNullString(cedu)
    t.RazonSocial = util.ValidarNullString(nomb)
    t.CuentaBancaria = util.ValidarNullString(cuen)
    t.Correo = util.ValidarNullString(corr)
    t.Ticket = util.ValidarNullString(tick)
    t.Serial = util.ValidarNullString(seri)
    t.Sistema = sist
    t.NombreSistema = util.ValidarNullString(sistema)
    t.MontoTicket = util.ValidarNullFloat64(montt)
    t.MontoSolicitado = util.ValidarNullFloat64(montos)
    t.Fecha = fecha

    lst = append(lst, t)
  }

  jSon, _ = json.Marshal(lst)
  return
}
