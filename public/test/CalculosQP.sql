

-- ##################################
-- GLOBAL 
-- ##################################
SELECT 
goid,g.obse,g.lote,g.parl,g.qued,g.part,g.calc,g.freq, venta,premio,comision FROM (
	SELECT  g.oid As goid,  SUM(vent) AS venta, SUM(prem) AS premio, SUM(comi) AS comision FROM (
		SELECT agen, fech, vent,prem,comi, sist from loteria
		UNION
		SELECT agen, fech, vent,prem,comi, sist from parley
	) AS A 
	JOIN zr_agencia zr ON A.agen=zr.codi
	JOIN agencia ON agencia.oid = zr.oida
	JOIN grupo g ON g.oid=zr.grupo
	JOIN sistema s ON s.oid=A.sist

	WHERE A.fech = '2017-01-03' 
	AND g.freq=4
	AND g.obse != '0'
	GROUP BY g.oid
	ORDER BY g.oid
) AS b
JOIN grupo g ON g.oid=B.goid


-- GLOBAL POR PROGRAMAS
SELECT goid,g.obse, 
	zrg.lote,zrg.parl,zrg.qued,zrg.part,zrg.calc,zrg.freq, 
	venta,premio,comision, b.sist,b.arch FROM (
	SELECT  g.oid As goid,  A.sist, s.arch, SUM(vent) AS venta, SUM(prem) AS premio, SUM(comi) AS comision FROM (
		SELECT agen, fech, vent,prem,comi, sist from loteria
		UNION
		SELECT agen, fech, vent,prem,comi, sist from parley
	) AS A 
	JOIN zr_agencia zr ON A.agen=zr.codi
	JOIN agencia ON agencia.oid = zr.oida
	JOIN grupo g ON g.oid=zr.grupo
	JOIN sistema s ON s.oid=A.sist

	WHERE A.fech BETWEEN '2017-01-01 00:00:00'::TIMESTAMP AND '2017-01-31 23:59:59'::TIMESTAMP
	AND g.freq=4
	AND g.obse != '0'
	
	
	GROUP BY g.oid, A.sist, s.arch
	ORDER BY g.oid
) AS b
JOIN grupo g ON g.oid=b.goid
LEFT JOIN zr_negociacion_grupo zrg ON zrg.oids=b.sist AND g.oid=zrg.oidg

--ALTER TABLE zr_negociacion_grupo
--ADD CONSTRAINT zr_negociacion_grupo_upkey UNIQUE(oidg, oids);


-- INSERT INTO zr_negociacion_grupo (oidg,oids,lote,parl,trip,term,qued,part,calc,freq)
-- select 1, oids,lote,parl,trip,term,qued,part,calc,freq from zr_negociacion_grupo WHERE oidg=8

--DELETE FROM zr_negociacion_grupo WHERE oidg=419



-- GLOBAL POR AGENCIAS
SELECT goid,g.obse, 
	zrg.lote,zrg.parl,zrg.qued,zrg.part,zrg.calc,zrg.freq, 
	venta,premio,comision, b.sist,b.arch FROM (
	SELECT  g.oid As goid,  A.agen, s.arch, SUM(vent) AS venta, SUM(prem) AS premio, SUM(comi) AS comision FROM (
		SELECT agen, fech, vent,prem,comi, sist from loteria
		UNION
		SELECT agen, fech, vent,prem,comi, sist from parley
	) AS A 
	JOIN zr_agencia zr ON A.agen=zr.codi
	JOIN agencia ON agencia.oid = zr.oida
	JOIN grupo g ON g.oid=zr.grupo
	JOIN sistema s ON s.oid=A.sist

	WHERE A.fech BETWEEN '2017-01-01 00:00:00'::TIMESTAMP AND '2017-01-31 23:59:59'::TIMESTAMP
	AND g.freq=4
	AND g.obse != '0'
	
	
	GROUP BY g.oid, A.agen, s.arch
	ORDER BY g.oid
) AS b
JOIN grupo g ON g.oid=b.goid
LEFT JOIN zr_negociacion_grupo zrg ON zrg.oids=b.sist AND g.oid=zrg.oidg