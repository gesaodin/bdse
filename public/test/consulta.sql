

/** POR AGENCIA **/

SELECT 
	lotepar.fech, SUM(lotepar.saldo) AS saldo, SUM(lotepar.vent) AS venta, 
	SUM(lotepar.prem) AS premio, SUM(lotepar.comi) AS comision

FROM agencia 
JOIN zr_agencia ON agencia.oid=zr_agencia.oida

JOIN (
	SELECT 
		arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria 
	
	UNION
	SELECT 
		arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
	
) AS lotepar ON zr_agencia.codi=lotepar.agen

WHERE agencia.obse='APMEMMPPAP00500' 
AND lotepar.fech BETWEEN '2017-01-01 00:00:00'::TIMESTAMP AND '2017-01-31 23:59:59'::TIMESTAMP
GROUP BY lotepar.fech
ORDER BY lotepar.fech



/** POR SISTEMA 
lotepar.fech, SUM(lotepar.saldo) AS saldo, SUM(lotepar.vent) AS venta, 
	SUM(lotepar.prem) AS premio, SUM(lotepar.comi) AS comision, lotepar.sist, 
	sistema.arch

**/
SELECT 
	lotepar.fech, SUM(lotepar.saldo) AS saldo, lotepar.sist, 
	sistema.obse, sistema.arch
FROM agencia 
JOIN zr_agencia ON agencia.oid=zr_agencia.oida

JOIN (
	SELECT 
		arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria 
	
	UNION
	SELECT 
		arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
	
) AS lotepar ON zr_agencia.codi=lotepar.agen
JOIN sistema ON lotepar.sist=sistema.oid
WHERE agencia.obse='APMEMMPPAP00500' 
AND lotepar.fech BETWEEN '2017-01-01 00:00:00'::TIMESTAMP AND '2017-01-02 23:59:59'::TIMESTAMP
GROUP BY lotepar.sist, lotepar.fech, sistema.arch, sistema.obse
ORDER BY sistema.arch







/** CONSULTA POR TAQUILLAS **/ 

SELECT 
	lotepar.agen, lotepar.fech, lotepar.saldo AS saldo, lotepar.vent AS venta, 
	lotepar.prem AS premio, lotepar.comi AS comision, lotepar.sist

FROM agencia 
JOIN zr_agencia ON agencia.oid=zr_agencia.oida

JOIN (
	SELECT 
		arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from loteria 
	
	UNION
	SELECT 
		arch, agen, fech, vent-prem-comi as saldo, vent, prem, comi, sist from parley
	
) AS lotepar ON zr_agencia.codi=lotepar.agen

WHERE agencia.obse='APMEMMPPAP00500' 
--AND lotepar.sist=9
AND lotepar.fech BETWEEN '2017-01-12 00:00:00'::TIMESTAMP AND '2017-01-12 23:59:59'::TIMESTAMP
--GROUP BY lotepar.sist, lotepar.fech
ORDER BY lotepar.agen

