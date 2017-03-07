package main

import (
	"fmt"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/util"
)

func init() {
	fmt.Println("")
	fmt.Println("Cargando Esquemas", sys.Version)
	fmt.Println("")
	sys.PostgresDBConexion()

}

func main() {
	var archivo = util.Archivo{}
	archivo.PostgreSQL = sys.PostgreSQL
	//archivo.LeerCodigosYCrearAgencias()
	archivo.LeerCodigosYCrearSaldos()
}
