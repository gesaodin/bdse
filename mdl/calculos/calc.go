package calculos

type Calculos struct {
	oid             int
	GananciaBruta   float64
	GananciaNeta    float64
	SaldoAculado    float64
	GastosMensuales float64
}

//Ejecutar Calculos para general datos
func (c *Calculos) Ejecutar() (jSon []byte, e error) {
	return
}

//Obtener Variables del sistema
func (c *Calculos) tenerVariables() (s string) {
	return
}
