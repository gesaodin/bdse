//subgrupo depende de un grupo
package subgrupo

import (
	"encoding/json"
	"strconv"

	"github.com/gesaodin/bdse/sys"
)

//SubGrupo Generacion de Grupos
type SubGrupo struct {
	ID               int          `json:"id,omitempty"` //Identificador
	Nombre           string       `json:"nombre,omitempty"`
	FechaNegociacion string       `json:"fecha,omitempty"`
	NumeroCuenta     string       `json:"cuenta,omitempty"`
	Triple           float64      `json:"triple,omitempty"`
	Terminal         float64      `json:"terminal,omitempty"`
	Queda            float64      `json:"queda,omitempty"`
	Participacion    float64      `json:"participacion,omitempty"`
	Observacion      string       `json:"observacion,omitempty"`
	Tipo             int          `json:"tipo,omitempty"`
	Frecuencia       int          `json:"frecuencia,omitempty"`  //1: Global 2:Individual
	Negociacion      int          `json:"negociacion,omitempty"` //1: Global 2:Individual
	Localizacion     Localizacion `json:"localizacion,omitempty"`
	Seguridad        Seguridad    `json:"seguridad,omitempty"`
}

//Localizacion Ubicacion Geografica
type Localizacion struct {
	IDParroquia int    `json:"idp,omitempty"`
	Casa        string `json:"casa,omitempty"`
	Direccion   string `json:"direccion,omitempty"`
	Telefono    string `json:"telefono,omitempty"`
	Celular     string `json:"celular,omitempty"`
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
func (g *SubGrupo) Registrar() (jSon []byte, err error) {
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

	s = `INSERT INTO grupo_localizacion (grupo,parro,casa,dire,cuen,tele,celu,obse) VALUES `
	s += `(` + strconv.Itoa(grupo) + `,` + parroquia + `,'` + g.Localizacion.Casa + `','`
	s += g.Localizacion.Direccion + `','` + g.NumeroCuenta + `','` + g.Localizacion.Telefono + `','` + g.Localizacion.Celular
	s += `','` + g.Observacion + `');`
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

//Consultar Grupo
func (g *SubGrupo) Consultar() (jSon []byte, err error) {
	return
}
