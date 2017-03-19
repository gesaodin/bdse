	SELECT oid, obse,
		vienen, saldo, entregado, recibido,ingreso, egreso, prestamo,cuota,
		saldo + vienen + (entregado - recibido) + (egreso - (ingreso+prestamo)) AS van, esta FROM (
			SELECT z.oid, z.obse,
			COALESCE(x.saldo,0) AS saldo,
			COALESCE(debe.monto,0) AS entregado,
			COALESCE(haber.monto,0) AS recibido,
			COALESCE(ingreso.monto,0) AS ingreso,
			COALESCE(egreso.monto,0) AS egreso,
			COALESCE(prestamo.monto,0) AS prestamo,
			COALESCE(prestamo.cuota,0) AS cuota,
			COALESCE(cobrosypagos.vien,0) AS vienen,
			COALESCE(cobrosypagoscierre.esta,0) AS esta
			FROM grupo g
			LEFT JOIN agencia z ON g.oid=z.grupo
			LEFT JOIN (
			SELECT agencia.oid, lotepar.fech, SUM(lotepar.saldo) AS saldo
			FROM agencia
			LEFT JOIN zr_agencia ON agencia.oid=zr_agencia.oida
			LEFT JOIN (
			SELECT agen, fech, vent-prem-comi as saldo from loteria
			UNION
			SELECT agen, fech, vent-prem-comi as saldo from parley
			) AS lotepar ON zr_agencia.codi=lotepar.agen

			WHERE  lotepar.fech BETWEEN '2016-12-31 00:00:00'::TIMESTAMP AND '2016-12-31 23:59:59'::TIMESTAMP 
			GROUP BY agencia.oid,lotepar.fech ) AS x ON x.oid=z.oid

			-- DEBE
			LEFT JOIN (
			SELECT agen, fapr, SUM(mont) AS monto FROM debe
			GROUP BY agen,fapr
			) AS debe ON
			debe.agen=z.obse AND debe.fapr='2016-12-31 00:00:00'::TIMESTAMP 

			-- HABER
			LEFT JOIN (
			SELECT agen, fapr, SUM(mont) AS monto FROM haber
			GROUP BY agen,fapr
			) AS haber ON
			haber.agen=z.obse  AND haber.fapr='2016-12-31 00:00:00'::TIMESTAMP 

			--INGRESO
			LEFT JOIN (
			SELECT agen, fapr, SUM(mont) AS monto FROM movimiento_ingreso
			GROUP BY agen,fapr
			)
			AS ingreso ON
			ingreso.agen=z.obse  AND ingreso.fapr='2016-12-31 00:00:00'::TIMESTAMP 

			-- EGRESO
			LEFT JOIN (
			SELECT agen, fapr, SUM(mont) AS monto FROM movimiento_egreso
			GROUP BY agen,fapr
			)
			AS egreso ON
			egreso.agen=z.obse  AND egreso.fapr='2016-12-31 00:00:00'::TIMESTAMP 

			-- PRESTAMOS
			LEFT JOIN (
			SELECT agen, fech, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
			GROUP BY agen,fech)
			AS prestamo ON
			prestamo.agen=z.obse AND prestamo.fech='2016-12-31 00:00:00'::TIMESTAMP 

			-- VIENEN
			LEFT JOIN cobrosypagos ON cobrosypagos.fech='2016-12-31 00:00:00'::TIMESTAMP 
			AND cobrosypagos.oida=z.oid

			-- CIERRE
			LEFT JOIN cobrosypagoscierre ON cobrosypagoscierre.fech='2016-12-31 00:00:00'::TIMESTAMP 
			WHERE  g.obse='0'
			ORDER BY z.obse
			) AS A


delete from cobrosypagos where fech='2017-01-02';
delete from cobrosypagos_grupo where fech='2017-01-02';
UPDATE cobrosypagos_grupo  SET sald=0, van=0, erec=0, movi=0 WHERE fech='2017-01-01';
delete from cobrosypagoscierre; --where fech='2017-01-02';
delete from cobrosypagoscierre_grupo; --where fech='2017-01-02';

