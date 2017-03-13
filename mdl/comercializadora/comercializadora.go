//compuesta por grupo, subgrupo, colector y agencia
package comercializadora

import (
	"strconv"

	"github.com/gesaodin/bdse/sys"
)

//Comercializadora Generacion de Grupos
type Comercializadora struct {
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
	Mensaje  string  `json:"msj,omitempty"`
	Tipo     int     `json:"tipo,omitempty"`
	Monto    float64 `json:"monto,omitempty"`
	Cantidad int     `json:"cantidad,omitempty"`
	Pgsql    string  `json:"pgsql,omitempty"`
}

//Registrar Salvar
func (c *Comercializadora) Registrar() (jSon []byte, err error) {

	return
}

//Consultar Grupo
func (c *Comercializadora) Consultar() (jSon []byte, err error) {
	return
}

//Cantidad de gupos asociados a una comercializadora
func (c *Comercializadora) Cantidad() (m Mensaje, err error) {
	//var m Mensaje
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
		//jSon, err = json.Marshal(m)
		//fmt.Println(err.Error())
		return
	}
	for sq.Next() {
		sq.Scan(&cantidad)
	}

	m.Tipo = 1
	m.Mensaje = strconv.Itoa(cantidad)
	m.Cantidad = cantidad
	//jSon, err = json.Marshal(m)
	return
}

//Gastos de la comercializadora movimientos de egresos
func (c *Comercializadora) Gastos() (m Mensaje, err error) {
	//var m Mensaje
	var gastos float64
	s := `
		SELECT COALESCE(SUM(mont), 0) as gastos FROM movimiento_egreso a
		WHERE a.comer=1 AND a.grupo=0 AND a.subgr=0
		AND a.colec=0;
	`
	sq, err := sys.PostgreSQL.Query(s)
	if err != nil {
		m.Mensaje = "Error: consulta de los gastos."
		m.Tipo = 2
		//jSon, err = json.Marshal(m)
		//fmt.Println(err.Error())
		return
	}
	for sq.Next() {
		sq.Scan(&gastos)
	}

	m.Tipo = 1
	m.Mensaje = strconv.FormatFloat(gastos, 'f', 6, 64)
	m.Monto = gastos

	//jSon, err = json.Marshal(m)
	return
}

//Depositos pendientes
func (c *Comercializadora) Depositos() (m Mensaje, err error) {
	var gastos float64
	s := `
		SELECT COALESCE(SUM(mont), 0) as deposito FROM haber a
		WHERE a.comer=1 AND a.grupo=0 AND a.subgr=0
		AND a.colec=0;
	`
	sq, err := sys.PostgreSQL.Query(s)
	if err != nil {
		m.Mensaje = "Error: consulta de los gastos."
		m.Tipo = 2
		//jSon, err = json.Marshal(m)
		//fmt.Println(err.Error())
		return
	}
	for sq.Next() {
		sq.Scan(&gastos)
	}

	m.Tipo = 1
	m.Mensaje = strconv.FormatFloat(gastos, 'f', 6, 64)
	m.Monto = gastos

	//jSon, err = json.Marshal(m)
	return
}
