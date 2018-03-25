package balance

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/util"
)

type Agencia struct {
	Descripcion   string  `json:"descripcion"`
	Codigo        string  `json:"codigo"`
	Fecha         string  `json:"fecha"`
	Venta         float64 `json:"venta"`
	Premio        float64 `json:"premio"`
	Comision      float64 `json:"comision"`
	Participacion float64 `json:"participacion"`
	Queda         float64 `json:"queda"`
	Calculo       float64 `json:"calculo"`
	NombreArchivo string  `json:"nombrearchivo"`
	Archivo       int     `json:"archivo"`
	Cantidad      int     `json:"cantidad"`
	Numero        int     `json:"numero"`
}

//ParticipacionSQL Consultando datos de participacion mayores a cero
func (a *Agencia) ParticipacionSQL(fecha string) string {
	return `SELECT VENTA.grupo, VENTA.oid, VENTA.obse, venta, premio, comision, 
	COALESCE(part,0) AS participacion,COALESCE(qued,0) AS queda,VENTA.programa,
	VENTA.archivo, calc, freq,
	((venta-premio-comision)*part)/100 AS calculo	
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
func (a *Agencia) CalcularParticipacionGlobal(fecha string) bool {

	// fmt.Println("Entrando en calculo")
	s := a.ParticipacionSQL(fecha)
	//fmt.Println(s)
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return false
	}
	for row.Next() {
		var grupo, oid, progr, arch int
		var vent, prem, comi, part, queda, calc, freq, mont sql.NullFloat64
		var obse string

		//Generar Movimientos de participación
		row.Scan(&grupo, &oid, &obse, &vent, &prem, &comi, &part, &queda, &progr, &arch, &calc, &freq, &mont)
		// venta := util.ValidarNullFloat64(vent)
		// premio := util.ValidarNullFloat64(prem)
		// comision := util.ValidarNullFloat64(comi)
		// participacion := util.ValidarNullFloat64(part)
		monto := util.ValidarNullFloat64(mont)
		smonto := strconv.FormatFloat(monto, 'f', 2, 64)
		s := insertMovimiento(grupo, oid, obse, fecha, smonto)
		_, err = sys.PostgreSQL.Query(s)
		if err != nil {
			return false
		}

	}
	return true
}

func insertMovimiento(grupo int, agencia int, desc string, fecha string, monto string) string {
	return `INSERT INTO movimiento_egreso 
		(comer,grupo,subgr,colec,agenc,agen,fech,fapr,fope,freg,tipo,cuen, oper,obse,mont) 
	VALUES 
		(
			1,` + strconv.Itoa(grupo) + `,0,0,` + strconv.Itoa(agencia) + `,
			'` + desc + `',
			'` + fecha + `'::DATE,'` + fecha + ` 00:00:00'::TIMESTAMP,
			'` + fecha + ` 00:00:00'::TIMESTAMP + '1 day',now(), 
			1, 0, '', 'PAGO POR PARTICIPACION - ` + fecha + `', ` + monto + `
		) `
}

//ValidarCajasSQL
func (a *Agencia) ValidarCajasSQL(fDesde string, fHasta string) string {
	return `SELECT agen, s.obse, vent, R.arch, ar.nomb,ar.fech,ar.cant,ar.tabl,ar.tipo  FROM (
		SELECT agen, vent, arch FROM zr_agencia AS T RIGHT JOIN (
		SELECT A.agen, SUM(vent) AS vent, A.arch  FROM (
			SELECT * FROM loteria
				UNION
			SELECT * FROM parley
				UNION
			SELECT * FROM figura
		) AS A
		GROUP BY A.agen, A.arch ) AS B ON T.codi=B.agen
		WHERE B.vent > 0 AND  T.codi IS NULL ) AS R
	JOIN archivo ar ON ar.oid=R.arch
	JOIN sistema s ON s.oid=ar.tipo
	WHERE ar.fech BETWEEN '` + fDesde + ` 00:00:00'::TIMESTAMP AND '` + fHasta + ` 23:59:59'::TIMESTAMP`
}

//ValidarCajas Validando cajas
func (a *Agencia) ValidarCajas(fDesde string, fHasta string) (jSon []byte, e error) {
	s := a.ValidarCajasSQL(fDesde, fHasta)
	// fmt.Println("SQL: ", s)
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}
	var lst []Agencia
	for row.Next() {
		var Agenc Agencia
		var vent sql.NullFloat64
		var agen, obse, nomb, fech string
		var cant, tabl, tipo, arch int
		row.Scan(&agen, &obse, &vent, &arch, &nomb, &fech, &cant, &tabl, &tipo)
		Agenc.Descripcion = obse
		Agenc.Codigo = agen
		Agenc.NombreArchivo = nomb
		Agenc.Venta = util.ValidarNullFloat64(vent)
		Agenc.Archivo = arch
		Agenc.Fecha = fech
		Agenc.Cantidad = cant
		lst = append(lst, Agenc)
	}
	jSon, e = json.Marshal(lst)
	return
}
