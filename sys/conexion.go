// PostgreSQL
// Es un Sistema de gestión de bases de datos relacional
// orientado a objetos y libre, publicado bajo la licencia PostgreSQL,
// similar a la BSD o la MIT.
package sys

import (
	"database/sql"

	"../util"
	_ "github.com/lib/pq"
)

// Funcion de Conexion a Postgres
func Conectar(c CadenaDeConexion) (db *sql.DB) {
	cadena := "user=" + c.Usuario + " dbname=" + c.BaseDatos + " password=" + c.Clave + " host=" + c.Host
	db, err := sql.Open("postgres", cadena)
	util.Error(err)
	return
}
