


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