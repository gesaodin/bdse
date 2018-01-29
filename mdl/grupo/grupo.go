//dependen de una comercializadora para existir
package grupo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/util"
)

//Grupo Generacion de Grupos
type Grupo struct {
	ID                 int          `json:"id,omitempty"`    //Identificador
	IDComercializadora int          `json:"comer,omitempty"` //Identificador
	Nombre             string       `json:"nombre,omitempty"`
	FechaNegociacion   string       `json:"fecha,omitempty"`
	NumeroCuenta       string       `json:"cuenta,omitempty"`
	Triple             float64      `json:"triple,omitempty"`
	Terminal           float64      `json:"terminal,omitempty"`
	Queda              float64      `json:"queda,omitempty"`
	Participacion      float64      `json:"participacion,omitempty"`
	Observacion        string       `json:"observacion,omitempty"`
	Frecuencia         int          `json:"frecuencia,omitempty"`  //1: Global 2:Individual
	Negociacion        int          `json:"negociacion,omitempty"` //1: Global 2:Individual
	Localizacion       Localizacion `json:"localizacion,omitempty"`
	Seguridad          Seguridad    `json:"seguridad,omitempty"`
	Saldos             Saldos       `json:"saldos,omitempty"`
	Tipo               int          `json:"tipo,omitempty"`
}

//Localizacion Ubicacion Geografica
type Localizacion struct {
	IDParroquia int    `json:"idp,omitempty"`
	Casa        string `json:"casa,omitempty"`
	Direccion   string `json:"direccion,omitempty"`
	Telefono    string `json:"telefono,omitempty"`
	Celular     string `json:"celular,omitempty"`
	Tipo        int    `json:"tipo,omitempty"` //1: Grupo 2: Subgrupo 3: Colector 4: Agencia
}

//Seguridad Validación de Acceso
type Seguridad struct {
	Usuario   string `json:"usuario,omitempty"`
	Correo    string `json:"correo,omitempty"`
	Clave     string `json:"clave,omitempty"`
	RClave    string `json:"rclave,omitempty"`
	Pregunta  int    `json:"pregunta,omitempty"`
	Respuesta string `json:"respuesta,omitempty"`
}

//Saldos Montos acumulados
type Saldos struct {
	Venta    float64 `json:"venta,omitempty"`
	Premio   float64 `json:"premio,omitempty"`
	Comision float64 `json:"comision,omitempty"`
	Saldo    float64 `json:"saldo,omitempty"`
}

//Mensaje del sistema
type Mensaje struct {
	Mensaje string `json:"msj,omitempty"`
	Tipo    int    `json:"tipo,omitempty"`
	Pgsql   string `json:"pgsql,omitempty"`
}

//Registrar Salvar
func (g *Grupo) Registrar() (jSon []byte, err error) {
	var grupo int
	var m Mensaje
	triple := strconv.FormatFloat(g.Triple, 'f', 6, 64)
	terminal := strconv.FormatFloat(g.Terminal, 'f', 6, 64)
	queda := strconv.FormatFloat(g.Queda, 'f', 6, 64)
	participacion := strconv.FormatFloat(g.Participacion, 'f', 6, 64)
	frecuencia := strconv.Itoa(g.Frecuencia)
	negociacion := strconv.Itoa(g.Negociacion)
	parroquia := strconv.Itoa(g.Localizacion.IDParroquia)
	s := `INSERT INTO grupo (comer,obse,resp,fneg,trip,term,qued,part,calc,freq,tipo) VALUES  `
	s += ` (1,'` + g.Nombre + `',1,'` + g.FechaNegociacion + `',` + triple + `,`
	s += terminal + `,` + queda + `,` + participacion + `,` + negociacion + `,` + frecuencia + `,0) RETURNING oid`

	sq, err := sys.PostgreSQL.Query(s)
	if err != nil {
		m.Mensaje = "Error: Grupo ya existe."
		m.Tipo = 2
		m.Pgsql = err.Error()
		jSon, err = json.Marshal(m)
		//fmt.Println(err.Error())
		return
	}
	for sq.Next() {
		sq.Scan(&grupo)
	}

	s = `INSERT INTO zr_gsca_localizacion (grupo,parro,casa,dire,cuen,tele,celu,obse,tipo,fech) VALUES `
	s += `(` + strconv.Itoa(grupo) + `,` + parroquia + `,'` + g.Localizacion.Casa + `','`
	s += g.Localizacion.Direccion + `','` + g.NumeroCuenta + `','` + g.Localizacion.Telefono + `','` + g.Localizacion.Celular
	s += `','` + g.Observacion + `',` + strconv.Itoa(g.Localizacion.Tipo) + `,now());`
	_, err = sys.PostgreSQL.Exec(s)
	if err != nil {
		m.Mensaje = "Error: ya existe la localización."
		m.Tipo = 2
		m.Pgsql = err.Error()
		jSon, err = json.Marshal(m)
		//fmt.Println(err.Error())
		return
	}

	s = `INSERT INTO usuario (nomb,ncom,corr,fech,esta,rol, toke) VALUES
				(
					'` + g.Seguridad.Usuario + `', 'Grupo Del Sistema','` + g.Seguridad.Correo + `',
					Now(), 1, 'Grupo', md5('` + g.Seguridad.Usuario + g.Seguridad.Clave + `')
				)`
	_, err = sys.PostgreSQL.Exec(s)
	if err != nil {
		m.Mensaje = "Error: Usuario ya existe."
		m.Tipo = 2
		m.Pgsql = err.Error()
		jSon, err = json.Marshal(m)
		//fmt.Println(err.Error())
		return
	}
	m.Tipo = 1
	m.Mensaje = "Proceso Exitoso"
	jSon, err = json.Marshal(m)
	return
}

/*
Consultar y convertir en una lista de objetos tipo Grupo.
Ofrece la localización que determina la ubicación geografica de cada
uno de los grupo así como sus reglas de negocio seguido de la seguridad
de acceso, correo y preguntas del sistema
*/
func (g *Grupo) Consultar() (LGrupo []Grupo, err error) {

	s := `
	-- Relación completa
SELECT g.obse, g.fneg, g.trip, g.term, g.qued, g.part,
		COALESCE(g.calc, 0) as calc,
		COALESCE(g.freq, 0) as freq,
		COALESCE(g.tipo, 0) as tipo,
		COALESCE(zr.parro, 0) as parro,
		zr.casa, zr.dire, zr.cuen, zr.tele, zr.celu, zr.obse, zr.fech,
		COALESCE(s.venta, 0) AS venta,
		COALESCE(s.premio, 0) AS premio,
		COALESCE(s.comision, 0) AS comision,
		COALESCE(s.saldo, 0) AS saldo
	FROM grupo g
	LEFT JOIN zr_gsca_localizacion zr ON g.oid=zr.grupo
	LEFT JOIN (
		SELECT
			g.oid,
			SUM(l.vent) AS venta,
			SUM(l.prem) AS premio,
			SUM(l.comi) AS comision,
			SUM(l.saldo) AS saldo
		FROM grupo g
		JOIN zr_agencia z ON g.oid=z.grupo
		JOIN (
			SELECT
				arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria
			UNION
			SELECT
				arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley

		) AS l ON z.codi=l.agen
		WHERE l.fech = (SELECT fech FROM cobrosypagoscierre ORDER BY fech desc LIMIT 1)
		GROUP BY g.oid
	) AS s ON s.oid = g.oid

	`
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		//m.Mensaje = "Error: TBL->Grupo"
		//fmt.Println(err.Error())
		return
	}
	for row.Next() {
		var gr Grupo
		var obse, fneg, casa, dire, cuent, tele, celu, fech, obser sql.NullString
		var trip, term, qued, part sql.NullFloat64
		var venta, premio, comision, saldo sql.NullFloat64
		var parr, calc, tipo, freq int

		e := row.Scan(
			&obse, &fneg, &trip, &term, &qued, &part, &calc,
			&freq, &tipo, &parr, &casa, &dire, &cuent, &tele,
			&celu, &obser, &fech,
			&venta, &premio, &comision, &saldo)

		if e != nil {
			fmt.Println(e.Error())
		}

		gr.Nombre = util.ValidarNullString(obse)
		gr.FechaNegociacion = util.ValidarNullString(fneg)
		gr.Triple = util.ValidarNullFloat64(trip)
		gr.Terminal = util.ValidarNullFloat64(term)
		gr.Queda = util.ValidarNullFloat64(qued)
		gr.Participacion = util.ValidarNullFloat64(part)
		gr.NumeroCuenta = util.ValidarNullString(cuent)
		gr.Frecuencia = freq
		gr.Tipo = tipo
		gr.Localizacion.IDParroquia = parr

		gr.Localizacion.Casa = util.ValidarNullString(casa)
		gr.Localizacion.Direccion = util.ValidarNullString(dire)
		gr.Localizacion.Telefono = util.ValidarNullString(tele)
		gr.Localizacion.Celular = util.ValidarNullString(celu)

		gr.Saldos.Venta = util.ValidarNullFloat64(venta)
		gr.Saldos.Premio = util.ValidarNullFloat64(premio)
		gr.Saldos.Comision = util.ValidarNullFloat64(comision)
		gr.Saldos.Saldo = util.ValidarNullFloat64(saldo)

		LGrupo = append(LGrupo, gr)

	}
	return
}

//Cantidad de gupos asociados a una comercializadora
func (g *Grupo) Cantidad() (jSon []byte, err error) {
	var m Mensaje
	var cantidad int
	s := `
		SELECT count(*) FROM comercializadora c
		JOIN grupo g on g.comer=c.oid
		WHERE c.oid=1
	`
	sq, err := sys.PostgreSQL.Query(s)
	if err != nil {
		m.Mensaje = "Error: consulta de grupo."
		m.Tipo = 2
		m.Pgsql = err.Error()
		jSon, err = json.Marshal(m)
		//fmt.Println(err.Error())
		return
	}
	for sq.Next() {
		sq.Scan(&cantidad)
	}

	m.Tipo = 1
	m.Mensaje = strconv.Itoa(cantidad)
	jSon, err = json.Marshal(m)
	return
}

//Gastos de la comercializadora egresos
func (g *Grupo) Gastos() (jSon []byte, err error) {
	var m Mensaje
	var gastos float64
	s := `
		SELECT SUM(mont) FROM movimiento_egreso a
		WHERE a.comer=1 AND a.grupo=0 AND a.subgr=0
		AND a.colec=0;
	`
	sq, err := sys.PostgreSQL.Query(s)
	if err != nil {
		m.Mensaje = "Error: consulta de los gastos."
		m.Tipo = 2
		jSon, err = json.Marshal(m)
		//fmt.Println(err.Error())
		return
	}
	for sq.Next() {
		sq.Scan(&gastos)
	}

	m.Tipo = 1
	m.Mensaje = strconv.FormatFloat(gastos, 'f', 6, 64)
	jSon, err = json.Marshal(m)
	return
}
