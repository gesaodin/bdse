// PostgreSQL
// Es un Sistema de gestión de bases de datos relacional
// orientado a objetos y libre, publicado bajo la licencia PostgreSQL,
// similar a la BSD o la MIT.
package sys

import (
	"database/sql"
	"fmt"

	mgo "gopkg.in/mgo.v2"

	"github.com/gesaodin/bdse/util"
	_ "github.com/lib/pq"
)

func MongoDBConexion() {
	MGOSession, Error = mgo.Dial("localhost:27000")
	fmt.Println("Cargando Conexión Con MongoDB...")
	util.Error(Error)
}

// Funcion de Conexion a Postgres
func PostgresDBConexion() {
	c := BaseDeDatos.CadenaDeConexion["postgres"]
	cadena := "user=" + c.Usuario + " dbname=" + c.Basedatos + " password=" + c.Clave + " host=" + c.Host
	PostgreSQL, _ = sql.Open("postgres", cadena)

	if PostgreSQL.Ping() != nil {
		fmt.Println("BDSE@Error: $ ", PostgreSQL.Ping())
	} else {
		fmt.Println("Conexión Establecida Con Postgres")
	}
	return
}
