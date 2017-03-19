-- ##################################
-- CONSULTA POR GRUPOS SALDOS DIARIOS
-- ##################################

SELECT 
	f.oid,f.obse,--f.lote,f.parl,f.qued,f.part,
	--f.calc,f.freq, 
	venta,premio,comision, comisioncal, (venta-premio-comisioncal) AS saldo,
	-- f.soid, --f.slote, f.sparl,f.squed, f.spart, 
	-- f.arch,f.fapr,
	entregado,recibido,ingreso,egreso,prestamo, cuota,
	vien,	
	(venta-premio-comisioncal) + vien + (entregado - recibido) + (ingreso-egreso+prestamo) AS van,
	COALESCE(cpc.esta,0) AS esta,cyp.fech
FROM 

	(
	SELECT 
		g.oid,g.obse,g.lote,g.parl,g.qued,g.part,g.calc,g.freq, 
		COALESCE(venta,0) AS venta,
		COALESCE(premio,0) AS premio,
		COALESCE(comision,0) AS comision, 
		CASE 
			WHEN g.lote > 0 then (venta * (g.lote + g.parl))/100
			WHEN g.parl > 0 then (venta * (g.lote + g.parl))/100
			WHEN zrg.lote > 0 then (venta * (zrg.lote + zrg.parl))/100
			WHEN zrg.parl > 0 then (venta * (zrg.lote + zrg.parl))/100
		ELSE 0
		END AS comisioncal, 
		b.soid, 
		zrg.lote AS slote, 
		zrg.parl AS sparl,
		zrg.qued AS squed, 
		zrg.part AS spart, 
		b.arch,
		egreso.fapr,
		COALESCE(debe.monto,0) AS entregado, 
		COALESCE(haber.monto,0) AS recibido,
		COALESCE(ingreso.monto,0) AS ingreso, 
		COALESCE(egreso.monto,0) AS egreso,
		COALESCE(prestamo.monto,0) AS prestamo,
		COALESCE(prestamo.cuota,0) AS cuota
		
	FROM 
	grupo g LEFT JOIN
	(
		SELECT  g.oid As goid,  SUM(vent) AS venta, SUM(prem) AS premio, SUM(comi) AS comision, s.oid as soid, s.arch FROM (
			SELECT agen, fech, vent,prem,comi, sist from loteria
			UNION
			SELECT agen, fech, vent,prem,comi, sist from parley
		) AS A 
		JOIN zr_agencia zr ON A.agen=zr.codi
		JOIN agencia ON agencia.oid = zr.oida
		JOIN grupo g ON g.oid=zr.grupo
		JOIN sistema s ON s.oid=A.sist

		WHERE A.fech = '2017-01-02' 
		AND g.freq=4
		--AND g.obse != '0'
		GROUP BY g.oid, s.oid,  s.arch
		ORDER BY g.oid
	) AS b ON g.oid=B.goid
	LEFT JOIN zr_negociacion_grupo zrg ON zrg.oids=b.soid AND g.oid=zrg.oidg

	-- DEBE
	LEFT JOIN (
			SELECT grupo, fapr, SUM(mont) AS monto FROM debe
			GROUP BY grupo,fapr
	) AS debe ON
	debe.grupo=g.oid AND debe.fapr='2017-01-02 00:00:00'::TIMESTAMP

	-- HABER
	LEFT JOIN (
		SELECT grupo, fapr, SUM(mont) AS monto FROM haber
		GROUP BY grupo,fapr
	) AS haber ON
	haber.grupo=g.oid AND haber.fapr='2017-01-02 00:00:00'::TIMESTAMP

	--INGRESO
	LEFT JOIN (
		SELECT grupo, fech,fapr, SUM(mont) AS monto FROM movimiento_ingreso
		GROUP BY grupo,fech,fapr
	)
	AS ingreso ON
	ingreso.grupo=g.oid AND ingreso.fapr='2017-01-02 00:00:00'::TIMESTAMP

	-- EGRESO
	LEFT JOIN (
		SELECT grupo, fech,fapr, SUM(mont) AS monto FROM movimiento_egreso
		GROUP BY grupo,fech,fapr
	)
	AS egreso ON
	egreso.grupo=g.oid AND egreso.fapr='2017-01-02 00:00:00'::TIMESTAMP

	-- PRESTAMOS
	LEFT JOIN (
		SELECT grupo, fapr, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
		GROUP BY grupo,fapr)
	AS prestamo ON
	prestamo.grupo=g.oid AND prestamo.fapr='2017-01-02 00:00:00'::TIMESTAMP


) AS f 	
LEFT JOIN cobrosypagos_grupo cyp ON cyp.fech='2017-01-02 00:00:00'::TIMESTAMP + '-24:00:00' AND cyp.oidg=f.oid
LEFT JOIN cobrosypagoscierre_grupo cpc ON cpc.fech='2017-01-02 00:00:00'::TIMESTAMP  + '-24:00:00'


--SELECT '2017-01-02 00:00:00'::TIMESTAMP + '-24:00:00';
/**
INSERT INTO cobrosypagos_grupo (oidg,fech,vien)
SELECT a.grupo,  c.fech, SUM(vien) FROM agencia a 
JOIN cobrosypagos c ON a.oid=c.oida
group by a.grupo, c.fech
order by a.grupo
**/