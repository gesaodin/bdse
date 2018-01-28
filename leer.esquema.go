package main

import (
	"fmt"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/util"
)

func leer() {
	fmt.Println("")
	fmt.Println("Cargando Esquemas", sys.Version)
	fmt.Println("")
	sys.PostgresDBConexion()

}

func main() {
	leer()
	var archivo = util.Archivo{}
	archivo.PostgreSQL = sys.PostgreSQL
	archivo.LeerCodigosYCrearSaldos()
	//archivo.LeerEntregados()
	//archivo.LeerEntregadosOficina()
	//archivo.LeerCodigosYCrearAgencias()
	//archivo.LeerCodigosYCrearSaldos()
}
