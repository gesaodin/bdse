//punto de ventas con taquillas las ventas generales
package agencia

import (
	"encoding/json"
	"strconv"

	"github.com/gesaodin/bdse/sys"
)

//Agencia Generacion de Grupos
type Agencia struct {
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
	Tipo        int    `json:"tipo,omitempty"`
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

//Caja La taquilla es el sitio donde se venden las entradas para acceder a
//un evento público, por ejemplo, al cine, al teatro o al estadio
type Caja struct {
	Oid    int    `json:"oid,omitempty"`
	Nombre string `json:"nombre,omitempty"`
	Fecha  string `json:"fecha,omitempty"`
}

//Sistema Programa de Ventas de Loterias: MATICLO, MORPHEUS, POS, PARLEY
type Sistema struct {
	Oid           int     `json:"oid,omitempty"`
	IDSistema     int     `json:"idsistema,omitempty"`
	Triple        float64 `json:"triple,omitempty"`
	Terminal      float64 `json:"terminal,omitempty"`
	Queda         float64 `json:"queda,omitempty"`
	Participacion float64 `json:"participacion,omitempty"`
	Fecha         string  `json:"fecha,omitempty"`
}

//Mensaje del sistema
type Mensaje struct {
	Mensaje string `json:"msj,omitempty"`
	Tipo    int    `json:"tipo,omitempty"`
	Pgsql   string `json:"pgsql,omitempty"`
}

//Registrar una agencia
func (a *Agencia) Registrar() (jSon []byte, err error) {
	var grupo int
	var m Mensaje
	triple := strconv.FormatFloat(a.Triple, 'f', 6, 64)
	terminal := strconv.FormatFloat(a.Terminal, 'f', 6, 64)
	queda := strconv.FormatFloat(a.Queda, 'f', 6, 64)
	participacion := strconv.FormatFloat(a.Participacion, 'f', 6, 64)
	frecuencia := strconv.Itoa(a.Frecuencia)
	negociacion := strconv.Itoa(a.Negociacion)
	parroquia := strconv.Itoa(a.Localizacion.IDParroquia)
	s := ` INSERT INTO grupo (comer,obse,resp,fneg,trip,term,qued,part,calc,freq,tipo) VALUES  `
	s += ` (1,'` + a.Nombre + `',1,'` + a.FechaNegociacion + `',` + triple + `,`
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
	s += `(` + strconv.Itoa(grupo) + `,` + parroquia + `,'` + a.Localizacion.Casa + `','`
	s += a.Localizacion.Direccion + `','` + a.NumeroCuenta + `','` + a.Localizacion.Telefono + `','` + a.Localizacion.Celular
	s += `','` + a.Observacion + `',` + strconv.Itoa(a.Localizacion.Tipo) + `,now());`
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
					'` + a.Seguridad.Usuario + `', 'Grupo Del Sistema','` + a.Seguridad.Correo + `',
					Now(), 1, 'Grupo', md5('` + a.Seguridad.Usuario + a.Seguridad.Clave + `')
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

//Consultar una Agencia
func (a *Agencia) Consultar() (jSon []byte, err error) {

	return
}

//Cantidad de gupos asociados a una comercializadora
func (a *Agencia) Cantidad() (jSon []byte, err error) {
	var m Mensaje
	var cantidad int
	s := `
		SELECT count(*) FROM comercializadora c
		JOIN agencia a ON a.comer=c.oid
		where c.oid=1 AND a.grupo=0 AND a.subgr=0
		AND a.colec=0
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

//Gastos de gupos asociados a una comercializadora
func (a *Agencia) Gastos() (jSon []byte, err error) {
	var m Mensaje
	var cantidad int
	s := `
		SELECT count(*) FROM comercializadora c
		JOIN agencia a ON a.comer=c.oid
		where c.oid=1 AND a.grupo=0 AND a.subgr=0
		AND a.colec=0
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
