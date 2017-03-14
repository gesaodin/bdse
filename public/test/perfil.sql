


-- Agencias
SELECT 
	a.obse, a.fneg,a.trip,a.term,a.qued,a.part,a.calc,a.freq,a.tipo,
	zr.parro,zr.casa, zr.dire, zr.cuen,zr.tele,zr.celu,
	zr.obse
FROM comercializadora c
	JOIN agencia a ON a.comer=c.oid
	LEFT JOIN zr_gsca_localizacion zr ON a.oid=zr.oida
WHERE c.oid=1 AND a.grupo=0 AND a.subgr=0
AND a.colec=0;


-- TOTAL DE GASTOS
SELECT SUM(mont) FROM movimiento_egreso a
WHERE a.comer=1 AND a.grupo=0 AND a.subgr=0
AND a.colec=0;

--DEPOSITOS PENDIENTES
SELECT SUM(mont) FROM haber a
WHERE a.comer=1 AND a.grupo=0 AND a.subgr=0
AND a.colec=0;

SELECT g.obse, g.fneg, g.trip, g.term, g.qued, g.part, g.calc, g.freq, g.tipo,
zr.parro, zr.casa, zr.dire, zr.cuen, zr.tele, zr.celu, zr.obse, zr.fech
FROM grupo g
LEFT JOIN zr_gsca_localizacion zr ON g.oid=zr.grupo

-- ventas, premios y comision: por grupos
SELECT 
	g.obse, 
	SUM(l.vent) AS venta,
	SUM(l.prem) AS premio,
	SUM(l.comi) AS comision,
	SUM(l.saldo) AS saldo
FROM grupo g 
JOIN zr_agencia z ON g.oid=z.grupo
JOIN (
	SELECT 
		arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria 
	UNION
	SELECT 
		arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
	
) AS l ON z.codi=l.agen
WHERE l.fech = (SELECT fech FROM cobrosypagoscierre ORDER BY fech desc LIMIT 1)
GROUP BY g.obse

-- Relación completa
SELECT g.obse, g.fneg, g.trip, g.term, g.qued, g.part,
			COALESCE(g.calc, 0) as calc,
			COALESCE(g.freq, 0) as freq,
			COALESCE(g.tipo, 0) as tipo,
			COALESCE(zr.parro, 0) as parro,
		 	zr.casa, zr.dire, zr.cuen, zr.tele, zr.celu, zr.obse, zr.fech,
		 	COALESCE(s.venta, 0) AS venta, 
		 	COALESCE(s.premio, 0) AS premio, 
		 	COALESCE(s.comision, 0) AS comision, 
		 	COALESCE(s.saldo, 0) AS saldo 
		FROM grupo g
		LEFT JOIN zr_gsca_localizacion zr ON g.oid=zr.grupo
		LEFT JOIN (
			SELECT 
				g.oid, 
				SUM(l.vent) AS venta,
				SUM(l.prem) AS premio,
				SUM(l.comi) AS comision,
				SUM(l.saldo) AS saldo
			FROM grupo g 
			JOIN zr_agencia z ON g.oid=z.grupo
			JOIN (
				SELECT 
					arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria 
				UNION
				SELECT 
					arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
				
			) AS l ON z.codi=l.agen
			WHERE l.fech = (SELECT fech FROM cobrosypagoscierre ORDER BY fech desc LIMIT 1)
			GROUP BY g.oid
		) AS s ON s.oid = g.oid











