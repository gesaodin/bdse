/** POR AGENCIA **/


SELECT saldo_agencia.fech, saldo_agencia.saldo,
	debe.monto AS entregado, haber.monto AS recibido,
	ingreso.monto AS ingreso, egreso.monto AS egreso,
	prestamo.monto AS prestamo,	
	cobrosypagos.vien AS vienen,
	cobrosypagos.van AS van

FROM (
SELECT agencia.obse, lotepar.fech, SUM(lotepar.saldo) AS saldo
FROM agencia 
JOIN zr_agencia ON agencia.oid=zr_agencia.oida
JOIN (
	SELECT arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria 	
	UNION
	SELECT arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley	
) AS lotepar ON zr_agencia.codi=lotepar.agen

WHERE agencia.obse='APMEMMPPAP00500' 
AND lotepar.fech BETWEEN '2017-01-01 00:00:00'::TIMESTAMP AND '2017-01-31 23:59:59'::TIMESTAMP
GROUP BY agencia.obse,lotepar.fech
) saldo_agencia
-- DEBE
LEFT JOIN (
		SELECT agen, fdep, SUM(mont) AS monto FROM debe
		GROUP BY agen,fdep
) AS debe ON
debe.agen=saldo_agencia.obse AND debe.fdep=saldo_agencia.fech
-- HABER
LEFT JOIN (
	SELECT agen, fdep, SUM(mont) AS monto FROM haber
	GROUP BY agen,fdep
) AS haber ON
haber.agen=saldo_agencia.obse  AND haber.fdep=saldo_agencia.fech

--INGRESO
LEFT JOIN (
	SELECT agen, fech, SUM(mont) AS monto FROM movimiento_ingreso
	GROUP BY agen,fech
)
AS ingreso ON
ingreso.agen=saldo_agencia.obse  AND ingreso.fech=saldo_agencia.fech
-- EGRESO
LEFT JOIN (
	SELECT agen, fech, SUM(mont) AS monto FROM movimiento_egreso
	GROUP BY agen,fech
)
AS egreso ON
egreso.agen=saldo_agencia.obse  AND egreso.fech=saldo_agencia.fech
-- PRESTAMOS
LEFT JOIN (
	SELECT agen, fech, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
	GROUP BY agen,fech)
AS prestamo ON
prestamo.agen=saldo_agencia.obse AND prestamo.fech=saldo_agencia.fech

-- VIENEN
LEFT JOIN cobrosypagos ON cobrosypagos.fech=saldo_agencia.fech

ORDER BY saldo_agencia.fech

select * from parley order by agen