package balance

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/util"
)

//Grupo
type Grupo struct{}

//ParticipacionSQLGlobalPP Consultando datos de participacion mayores a cero Por Programas
func (g *Grupo) ParticipacionSQLGlobalPP(fecha string) string {
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

//CalcularParticipacionGlobal Función para ejecutar calculos de participaciones por agencia
func (g *Grupo) CalcularParticipacionGlobalPP(fecha string) bool {

	// fmt.Println("Entrando en calculo")
	s := g.ParticipacionSQLGlobalPP(fecha)
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
		tabla := "movimiento_egreso"
		if monto < 0 {
			tabla = "movimiento_ingreso"
			monto = monto * -1
		}
		smonto := strconv.FormatFloat(monto, 'f', 2, 64)

		s := insertMovimientoG(oidg, obse, fecha, smonto, tabla, "PARTICIPACION ")
		_, err = sys.PostgreSQL.Query(s)
		if err != nil {
			return false
		}

	}
	return true
}

func insertMovimientoG(grupo int, desc string, fecha string, monto string, tabla string, modelo string) string {
	return `INSERT INTO  ` + tabla + `
		(comer,grupo,subgr,colec,agenc,agen,fech,fapr,fope,freg,tipo,cuen, oper,obse,mont) 
	VALUES 
		(
			1,` + strconv.Itoa(grupo) + `,0,0,0,
			'` + desc + `',
			'` + fecha + `'::DATE,'` + fecha + ` 00:00:00'::TIMESTAMP,
			'` + fecha + ` 00:00:00'::TIMESTAMP + '1 day',now(), 
			1, 0, '', 'PAGO POR ` + modelo + ` - ` + fecha + `', ` + monto + `
		) `
}

//ParticipacionSQLGlobalG Consultando datos de participacion mayores a cero
func (g *Grupo) ParticipacionSQLGlobalG(fecha string) string {
	return `	
		SELECT grupo,obse,venta,premio,comision,part,calc, ((venta-premio-comision)*part)/100 AS calculo FROM (
			SELECT  zr.grupo, SUM(vent) AS venta, SUM(prem) AS premio, 
				SUM(comi) AS comision FROM (
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
			GROUP BY zr.grupo
		) AS G
		JOIN grupo gr ON G.grupo=gr.oid
		WHERE calc=2 AND part>0	
		ORDER BY gr.oid`
}

//CalcularParticipacionGlobalG Función para ejecutar calculos de participaciones por agencia
func (g *Grupo) CalcularParticipacionGlobalG(fecha string) bool {
	s := g.ParticipacionSQLGlobalG(fecha)
	//fmt.Println(s)
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return false
	}
	for row.Next() {
		var oidg int
		var obse string
		var venta, premio, comision, part, calc, mont sql.NullFloat64
		row.Scan(&oidg, &obse, &venta, &premio, &comision, &part, &calc, &mont)
		monto := util.ValidarNullFloat64(mont)
		tabla := "movimiento_egreso"
		if monto < 0 {
			tabla = "movimiento_ingreso"
			monto = monto * -1
		}
		smonto := strconv.FormatFloat(monto, 'f', 2, 64)

		s := insertMovimientoG(oidg, obse, fecha, smonto, tabla, "PARTICIPACION ")
		_, err = sys.PostgreSQL.Query(s)
		if err != nil {
			return false
		}
	}
	return true
}

/**************************************
**  Casos vistos desde el Grupo
***************************************/

//CQGlobal Calcular Queda Global
//Caso 1
func (g *Grupo) CQGlobal(fecha string, freq int) {
	var Queda Queda
	Queda.Fecha = validarFrecuencia(fecha, freq)
	Queda.Tipo = strconv.Itoa(freq)
	s := Queda.AGlobal()
	fmt.Println(s)
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}
	for row.Next() {
		var grupo int
		var obse, qued, freq string
		var saldo sql.NullFloat64
		row.Scan(&grupo, &obse, &qued, &freq, &saldo)
		// monto := util.ValidarNullFloat64(saldo)
		// smonto := strconv.FormatFloat(monto, 'f', 2, 64)
		// s := insertMovimiento(grupo, oid, obse, fecha, smonto, "QUEDA")
		// _, err = sys.PostgreSQL.Query(s)
		// if err != nil {
		// 	return
		// }
	}
}
