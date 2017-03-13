//ejecución de formulas para ganancias y perdidas
package calculos

//Calculos Para Ganancias
type Calculos struct {
	oid             int     //Identificadors
	GananciaBruta   float64 //Ganancias Brutas
	GananciaNeta    float64 // Ganacias Netas
	SaldoAculado    float64 //Saldo Actual
	GastosMensuales float64 //Gastos Mensuales
}

//Ejecutar Calculos para general datos
func (c *Calculos) Ejecutar() (jSon []byte, e error) {
	return
}

//Obtener Variables del sistema
func (c *Calculos) tenerVariables() (s string) {
	return
}
