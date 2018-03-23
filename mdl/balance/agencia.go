package balance

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/util"
)

type Agencia struct{}

//Participacion Consultando datos de participacion mayores a cero
func (a *Agencia) Participacion(fecha string) string {
	return `SELECT VENTA.grupo, VENTA.oid, VENTA.obse, venta, premio, comision, 
	COALESCE(part,0) AS participacion,COALESCE(qued,0) AS queda,VENTA.programa,VENTA.archivo
	FROM zr_negociacion_agencia AS AGN
	RIGHT JOIN (
		SELECT  zr.grupo, agencia.oid, agencia.obse, SUM(vent) AS venta, SUM(prem) AS premio, SUM(comi) AS comision, s.oid as programa, s.arch AS archivo FROM (
			SELECT agen, fech, vent,prem,comi, sist from loteria
			UNION
			SELECT agen, fech, vent,prem,comi, sist from parley
			UNION
			SELECT agen, fech, vent,prem,comi, sist from figura
		) AS A
		JOIN zr_agencia zr ON A.agen=zr.codi
		JOIN agencia ON agencia.oid = zr.oida
		JOIN sistema s ON s.oid=A.sist
		WHERE A.fech = '` + fecha + `'
		GROUP BY zr.grupo, agencia.oid, agencia.obse, s.oid,  s.arch
		ORDER BY agencia.oid
	) AS VENTA 
		ON AGN.oida=VENTA.oid AND AGN.oids=VENTA.programa
	WHERE AGN.part > 0`
}

//CalcularParticipacion Función para ejecutar calculos de participaciones por agencia
func (a *Agencia) CalcularParticipacion(fecha string) {
	fmt.Println("Entrando en calculo")
	s := a.Participacion(fecha)
	fmt.Println(s)
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}
	for row.Next() {
		var grupo, oid, progr, arch int
		var vent, prem, comi, part, queda sql.NullFloat64
		var obse string
		var monto float64

		//Generar Movimientos de participación
		row.Scan(&grupo, &oid, &obse, &vent, &prem, &comi, &part, &queda, &progr, &arch)
		venta := util.ValidarNullFloat64(vent)
		premio := util.ValidarNullFloat64(prem)
		comision := util.ValidarNullFloat64(comi)
		participacion := util.ValidarNullFloat64(part)

		fmt.Println(venta, premio, comision, participacion)

		monto = ((venta - premio - comision) * participacion) / 100
		smonto := strconv.FormatFloat(monto, 'f', 2, 64)
		s := insertMovimiento(grupo, oid, obse, fecha, smonto)
		fmt.Println(s)

	}
}

func insertMovimiento(grupo int, agencia int, desc string, fecha string, monto string) string {
	return `INSERT INTO movimiento_egreso 
		(comer,grupo,subgr,colec,agenc,fech,fapr,freg,tipo,cuen, oper,obse,mont) 
	VALUES 
		(
			1,` + strconv.Itoa(grupo) + `,0,0,` + strconv.Itoa(agencia) + `,
			'` + fecha + `'::DATE,'` + fecha + ` 00:00:00'::TIMESTAMP + '1 day',now(), 
			1, 0, '', 'PAGO POR PARTICIPACION - ` + fecha + `', ` + monto + `
		) `
}
