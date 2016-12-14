// El consumidor de recursos establece la necesidad de uso e iteraciones
// a las cuales el usuario se enfrenta a diario. Recurrencia
package logs

import "time"

const (
	Irregular = 0
	Regular   = 1
)

type ValorEsperado struct{}

type Comportamiento struct {
	tipo   int
	Accion string
	ValorEsperado
}

type Anomalia struct {
	Comportamiento
	funcion string
	tiempo  time.Time
}

func (r *Anomalia) Agregar() {

}

func (r *Anomalia) Notificar() {

}
