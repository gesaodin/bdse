//modelos de objets del sistema
package mdl

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gesaodin/bdse/sys"

	"gopkg.in/mgo.v2/bson"
)

//NombreCompleto Nombres y Apellidos
type NombreCompleto struct {
	Nombre   string `json:"nombre" bson:"nombre"`
	Apellido string `json:"apellido" bson:"apellido"`
}

//Persona El significado actual de persona tiene su origen en las controversias
// cristológicas de los siglos IV y V. En el transcurso del debate entre
// las diferentes escuelas teológicas, se desarrollaron conceptos hasta
// entonces no conocidos.
type Persona struct {
	ID                bson.ObjectId          `json:"id" bson:"_id"`
	Cedula            string                 `json:"cedula" bson:"cedula"`
	Pasaporte         string                 `json:"pasaporte" bson:"pasaporte"`
	RIF               string                 `json:"rif" bson:"rif"`
	NombreCompleto    interface{}            `json:"nombrecompleto"`
	Nacionalidad      string                 `json:"nacionalidad" bson:"nacionalidad"`
	Sexo              string                 `json:"sexo" bson:"sexo"`
	FechaDeNacimiento string                 `json:"fechadenacimiento" bson:"fechadenacimiento"`
	FechaDeCreacion   time.Time              `json:"fechadecreacion" bson:"fechadecreacion"`
	Direccion         map[string]interface{} `json:"direccion" bson:"direccion"`
	//Telefonos         []Telefono `json:"telefonos" bson:"telefonos"`
}

//Telefono (XXXX) -XXX-XXXX
type Telefono struct {
	Tipo         string `json:"tipo" bson:"tipo"`
	CodigoDeArea string `json:"codigodearea" bson:"codigodearea"`
	Numero       string `json:"numero" bson:"numero"`
}

//Direccion localización o ubicación
type Direccion struct {
	Estado    string
	Ciudad    string
	Municipio string
	Parroquia string
	Sector    string
	Detalle   string
}

//DBPersona Objeto en Base de Datos con campos null
type DBPersona struct {
	Cedula          string
	Nombre          sql.NullString
	Apellido        sql.NullString
	Nacionalidad    sql.NullString
	Sexo            sql.NullString
	FechaNacimiento sql.NullString
	//Telefonos	[]Telefono
	//Direcciones	[]Direccion
}

//ConsultarMGO una persona mediante el metodo de MongoDB
func (p *Persona) ConsultarMGO(cedula string) (err error) {
	c := sys.MGOSession.DB("bdse").C("persona")
	err = c.Find(bson.M{"cedula": cedula}).One(&p)
	return
}

//ListarMGO Listado General
func (p *Persona) ListarMGO(cedula string) (lst []Persona, err error) {
	c := sys.MGOSession.DB("bdse").C("persona")
	err = c.Find(bson.M{}).All(&lst)
	return
}

//SalvarMGO Guardar
func (p *Persona) SalvarMGO() (err error) {
	c := sys.MGOSession.DB("bdse").C("persona")
	//fmt.Println(p)
	err = c.Insert(p)
	return
}

//ActualizarMGO Actualizar
func (p *Persona) ActualizarMGO(persona map[string]interface{}) (err error) {
	c := sys.MGOSession.DB("bdse").C("persona")
	err = c.Update(bson.M{"cedula": persona["cedula"]}, bson.M{"$set": persona})

	return
}

//ListarPostgreSQL Listar Personas en Postgres
func (p *Persona) ListarPostgreSQL() {
	fmt.Println("Entrando PostgreSQL")
	rows, error := sys.PostgreSQL.Query("SELECT cedula FROM beneficiario LIMIT 10")
	if error != nil {
		panic(error)
	}
	defer rows.Close()
	for rows.Next() {
		var cedula string
		if err := rows.Scan(&cedula); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s \n", cedula)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
