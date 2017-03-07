
SELECT oid, obse, 
vienen, saldo, entregado, recibido,ingreso, egreso, prestamo,
saldo + vienen + (entregado - recibido) + (egreso - (ingreso+prestamo)) AS van FROM (
SELECT z.oid, obse,
	COALESCE(x.saldo,0) AS saldo, 
	COALESCE(debe.monto,0) AS entregado,
	COALESCE(haber.monto,0) AS recibido,
	COALESCE(ingreso.monto,0) AS ingreso, 
	COALESCE(egreso.monto,0) AS egreso,
	COALESCE(prestamo.monto,0) AS prestamo,
	COALESCE(cobrosypagos.vien,0) AS vienen
FROM agencia AS z
LEFT JOIN (
	SELECT agencia.oid, lotepar.fech, SUM(lotepar.saldo) AS saldo
		FROM agencia
	LEFT JOIN zr_agencia ON agencia.oid=zr_agencia.oida
	LEFT JOIN (
		SELECT agen, fech, vent-prem-comi as saldo from loteria
		UNION
		SELECT agen, fech, vent-prem-comi as saldo from parley
	) AS lotepar ON zr_agencia.codi=lotepar.agen

	WHERE lotepar.fech BETWEEN '2017-01-01 00:00:00'::TIMESTAMP AND '2017-01-01 23:59:59'::TIMESTAMP
	
	GROUP BY agencia.oid,lotepar.fech ) AS x ON x.oid=z.oid

-- DEBE
LEFT JOIN (
	SELECT agen, fapr, SUM(mont) AS monto FROM debe
	GROUP BY agen,fapr
) AS debe ON
debe.agen=z.obse AND debe.fapr=x.fech

-- HABER
LEFT JOIN (
	SELECT agen, fapr, SUM(mont) AS monto FROM haber
	GROUP BY agen,fapr
) AS haber ON
haber.agen=z.obse  AND haber.fapr=x.fech

--INGRESO
LEFT JOIN (
	SELECT agen, fech, SUM(mont) AS monto FROM movimiento_ingreso
	GROUP BY agen,fech
)
AS ingreso ON
ingreso.agen=z.obse  AND ingreso.fech=x.fech

-- EGRESO
LEFT JOIN (
	SELECT agen, fech, SUM(mont) AS monto FROM movimiento_egreso
	GROUP BY agen,fech
)
AS egreso ON
egreso.agen=z.obse  AND egreso.fech=x.fech

-- PRESTAMOS
LEFT JOIN (
	SELECT agen, fech, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
	GROUP BY agen,fech)
AS prestamo ON
prestamo.agen=z.obse AND prestamo.fech=x.fech

-- VIENEN
LEFT JOIN cobrosypagos ON cobrosypagos.fech=x.fech
AND cobrosypagos.agen=z.obse
ORDER BY z.obse
) AS A