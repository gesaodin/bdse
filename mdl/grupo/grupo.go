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
	s := ` INSERT INTO grupo (comer,obse,resp,fneg,trip,term,qued,part,calc,freq,tipo) VALUES  `
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
		SELECT g.obse, g.fneg, g.trip, g.term, g.qued, g.part, g.calc, g.freq, g.tipo,
		zr.parro, zr.casa, zr.dire, zr.cuen, zr.tele, zr.celu, zr.obse, zr.fech
		FROM grupo g
		LEFT JOIN zr_gsca_localizacion zr ON g.oid=zr.grupo
	`
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		//m.Mensaje = "Error: TBL->Grupo"
		//fmt.Println(err.Error())
		return
	}
	for row.Next() {
		var gr Grupo
		var obse, fneg, casa, dire, cuent, tele, celu, fech, obser string
		var trip, term, qued, part sql.NullFloat64
		var parr, calc, tipo, freq int

		e := row.Scan(
			&obse, &fneg, &trip, &term, &qued, &part, &calc,
			&freq, &tipo, &parr, &casa, &dire, &cuent, &tele,
			&celu, &obser, &fech)

		if e != nil {
			fmt.Println(err.Error())
		}

		gr.FechaNegociacion = fneg
		gr.Triple = util.ValidarNullFloat64(trip)
		gr.Terminal = util.ValidarNullFloat64(term)
		gr.Queda = util.ValidarNullFloat64(qued)
		gr.Participacion = util.ValidarNullFloat64(part)
		gr.NumeroCuenta = cuent
		gr.Frecuencia = freq
		gr.Tipo = tipo
		gr.Localizacion.IDParroquia = parr

		gr.Localizacion.Casa = casa
		gr.Localizacion.Direccion = dire
		gr.Localizacion.Telefono = tele
		gr.Localizacion.Celular = celu

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
