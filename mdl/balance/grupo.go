package balance

import (
	"database/sql"
	"strconv"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/util"
)

type Grupo struct{}

//ParticipacionSQL Consultando datos de participacion mayores a cero
func (g *Grupo) ParticipacionSQL(fecha string) string {
	return `	
		SELECT  oidg, obse,part,calc,freq, oids,venta,premio,comision, 
		((venta-premio-comision)*part)/100 AS calculo	
		FROM (
		SELECT zrg.oidg, gr.obse,zrg.part,zrg.calc,zrg.freq, zrg.oids FROM grupo gr
		JOIN zr_negociacion_grupo zrg ON gr.oid=zrg.oidg
		WHERE zrg.part>0
		ORDER BY zrg.oidg,zrg.oids
		) AS A
		JOIN (
			SELECT grupo, SUM(venta) venta, SUM(premio) premio, SUM(comision) comision, programa, archivo FROM (
					SELECT  zr.grupo, agencia.oid, agencia.obse, SUM(vent) AS venta, SUM(prem) AS premio, 
						SUM(comi) AS comision, s.oid as programa, s.arch AS archivo FROM (
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
			) AS G
			
			GROUP BY grupo, programa, archivo
			ORDER BY grupo
		) AS rgr
		ON A.oidg = rgr.grupo AND A.oids=rgr.programa`
}

//CalcularParticipacionIndividual Función para ejecutar calculos de participaciones por agencia
func (g *Grupo) CalcularParticipacionIndividual(fecha string) bool {

	// fmt.Println("Entrando en calculo")
	s := g.ParticipacionSQL(fecha)
	//fmt.Println(s)
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return false
	}
	for row.Next() {
		var oidg, progr int
		var part, calc, freq, vent, prem, comi, mont sql.NullFloat64
		var obse string

		//Generar Movimientos de participación
		row.Scan(&oidg, &obse, &part, &calc, &freq, &progr, &vent, &prem, &comi, &mont)
		// venta := util.ValidarNullFloat64(vent)
		// premio := util.ValidarNullFloat64(prem)
		// comision := util.ValidarNullFloat64(comi)
		// participacion := util.ValidarNullFloat64(part)
		monto := util.ValidarNullFloat64(mont)
		smonto := strconv.FormatFloat(monto, 'f', 2, 64)
		s := insertMovimientoG(oidg, obse, fecha, smonto)
		_, err = sys.PostgreSQL.Query(s)
		if err != nil {
			return false
		}

	}
	return true
}

func insertMovimientoG(grupo int, desc string, fecha string, monto string) string {
	return `INSERT INTO movimiento_egreso 
		(comer,grupo,subgr,colec,agenc,agen,fech,fapr,fope,freg,tipo,cuen, oper,obse,mont) 
	VALUES 
		(
			1,` + strconv.Itoa(grupo) + `,0,0,0,
			'` + desc + `',
			'` + fecha + `'::DATE,'` + fecha + ` 00:00:00'::TIMESTAMP,
			'` + fecha + ` 00:00:00'::TIMESTAMP + '1 day',now(), 
			1, 0, '', 'PAGO POR PARTICIPACION - ` + fecha + `', ` + monto + `
		) `
}
