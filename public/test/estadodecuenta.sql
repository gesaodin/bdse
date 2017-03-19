SELECT * FROM grupo g
JOIN cobrosypagos_grupo cg ON g.oid=cg.oidg
JOIN zr_agencia zr ON g.oid=zr.grupo
LEFT JOIN
	(
		SELECT agen, fech, vent,prem,comi, sist from loteria
		UNION
		SELECT agen, fech, vent,prem,comi, sist from parley
	) AS A ON A.agen=zr.codi

where g.oid=408
AND A.fech='2017-01-03'


-- ############################
-- ESTADO DE CUENTA DEL GRUPO
-- ############################

SELECT fech,vien, sald, movi, erec, van FROM grupo g
JOIN cobrosypagos_grupo cg ON g.oid=cg.oidg
where g.oid=408
AND van!=0
ORDER BY cg.fech

SELECT fech,vien, sald, movi, erec, van FROM grupo g
JOIN cobrosypagos_grupo cg ON g.oid=cg.oidg
where g.oid=408 AND cg.fech BETWEEN '2017-01-01 00:00:00'::TIMESTAMP AND '2017-01-31 23:59:59'::TIMESTAMP
AND van!=0
