package main

import (
	"fmt"

	"github.com/gesaodin/bdse/sys"
)

func init() {
	if sys.Postgres {
		fmt.Println("Activando Posgres")
	}
}

func main() {
	fmt.Println("Cargando Bus de Servicio Empresarial")
}
