//configuraciones del sistema
package sys

import (
	"database/sql"
	"fmt"

	mgo "gopkg.in/mgo.v2"

	"github.com/gesaodin/bdse/util"
	_ "github.com/lib/pq"
)

//MongoDBConexion Conexion a Mongo DB
func MongoDBConexion() {
	MGOSession, Error = mgo.Dial("localhost:27000")
	fmt.Println("Cargando Conexión Con MongoDB...")
	util.Error(Error)
}

//PostgresDBConexion Funcion de Conexion a Postgres
func PostgresDBConexion() {
	c := BaseDeDatos.CadenaDeConexion["postgres"]
	cadena := "user=" + c.Usuario + " dbname=" + c.Basedatos + " password=" + c.Clave + " host=" + c.Host + " sslmode=disable"
	PostgreSQL, _ = sql.Open("postgres", cadena)

	if PostgreSQL.Ping() != nil {
		fmt.Println("BDSE@Error: $ ", PostgreSQL.Ping())
	} else {
		fmt.Println("Conexión Establecida Con Postgres")
	}
	return
}
