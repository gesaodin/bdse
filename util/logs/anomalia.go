// El consumidor de recursos establece la necesidad de uso e iteraciones
// a las cuales el usuario se enfrenta a diario. Recurrencia
package logs

import (
	"fmt"
	"time"
)

const (
	Irregular = 0
	Regular   = 1
)

type ValorEsperado struct {
	Parametro string
	Resultado interface{}
}
type Comportamiento struct {
	tipo   int
	Accion string
}

type Anomalia struct {
	Comportamiento
	funcion string
	tiempo  time.Time
}

func (r *Anomalia) Agregar() {

}

func (r *Anomalia) Notificar(v ValorEsperado) {

	switch r := v.Resultado.(type) {
	case string:
		fmt.Println("String: ", r)
	case int, int64:
		fmt.Println("Entero: ", r)
	case float32, float64:
		fmt.Println("Float: ", r)
	default:

		fmt.Println("Valor Inesperado", r)
	}
}
