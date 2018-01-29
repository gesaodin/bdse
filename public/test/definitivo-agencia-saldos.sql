	SELECT cpc.oid, cpc.fech,
			cyp.vien,
			saldo_agencia.saldo,
			debe.monto AS entregado,
			haber.monto AS recibido,
			ingreso.monto AS ingreso,
			egreso.monto AS egreso,
			prestamo.monto AS prestamo,
			prestamo.cuota AS cuota,
			COALESCE(saldo,0) + COALESCE(cyp.vien,0) +
			(
				COALESCE(debe.monto,0) - COALESCE(haber.monto,0)) +
				(
					COALESCE(egreso.monto,0) - (COALESCE(ingreso.monto,0)+COALESCE(prestamo.monto,0))
				) AS van,
			cpc.esta
			FROM cobrosypagoscierre cpc
			

			--) AS f
			--ON cpc.fech=f.fech

			-- VIENEN
			LEFT JOIN cobrosypagos cyp ON cyp.fech=cpc.fech
			INNER JOIN agencia ON cyp.oida=agencia.oid

			
			-- DEBE
			LEFT JOIN (
					SELECT agen, fapr AS fapr, SUM(mont) AS monto FROM debe
					GROUP BY agen,fapr
			) AS debe ON
			debe.agen=agencia.obse AND debe.fapr=cyp.fech

			-- HABER
			LEFT JOIN (
				SELECT agen, fapr, SUM(mont) AS monto FROM haber
				GROUP BY agen,fapr
			) AS haber ON
			haber.agen=agencia.obse  AND haber.fapr=cyp.fech

			--INGRESO
			LEFT JOIN (
				SELECT agen, fech,fapr, SUM(mont) AS monto FROM movimiento_ingreso
				GROUP BY agen,fech,fapr
			)
			AS ingreso ON
			ingreso.agen=agencia.obse  AND ingreso.fapr=cyp.fech

			-- EGRESO
			LEFT JOIN (
				SELECT agen, fech,fapr, SUM(mont) AS monto FROM movimiento_egreso
				GROUP BY agen,fech,fapr
			)
			AS egreso ON
			egreso.agen=agencia.obse  AND egreso.fapr=cyp.fech

			-- PRESTAMOS
			LEFT JOIN (
				SELECT agen, fech, SUM(mont) AS monto, SUM(mcuo) AS cuota FROM movimiento_prestamo
				GROUP BY agen,fech)
			AS prestamo ON
			prestamo.agen=agencia.obse AND prestamo.fech=cyp.fech


			LEFT JOIN (
				--SELECT saldo_agencia.oid, saldo_agencia.fech, saldo_agencia.saldo,
				--	debe.monto AS entregado, haber.monto AS recibido,
				--	ingreso.monto AS ingreso, egreso.monto AS egreso,
				--	prestamo.monto AS prestamo,prestamo.cuota AS cuota
				--FROM (
				SELECT agencia.oid,agencia.obse, lotepar.fech, SUM(lotepar.saldo) AS saldo
				FROM agencia
				JOIN zr_agencia ON agencia.oid=zr_agencia.oida
				JOIN (
					SELECT arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria
					UNION
					SELECT arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
				) AS lotepar ON zr_agencia.codi=lotepar.agen

				WHERE agencia.obse='APMEMMPPCD001'  
				AND lotepar.fech BETWEEN '2017-01-01 00:00:00'::TIMESTAMP 
				AND '2017-01-31 23:59:59'::TIMESTAMP
				GROUP BY agencia.oid,agencia.obse,lotepar.fech
			) AS saldo_agencia ON saldo_agencia.obse=agencia.obse AND saldo_agencia.fech=cyp.fech

			
			WHERE agencia.obse='APMEMMPPCD001' AND cyp.oida=agencia.oid
			ORDER BY fech