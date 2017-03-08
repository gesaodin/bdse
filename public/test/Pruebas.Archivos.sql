SELECT DISTINCT agen, sistema.obse FROM loteria
join sistema on sistema.oid=loteria.sist
UNION
SELECT DISTINCT agen, sistema.obse FROM parley
join sistema on sistema.oid=parley.sist
ORDER BY agen

select * from loteria where agen='APMEBAPPRH00300(603)'
select * from loteria where agen='APMEMMPPRO00100(31)'

select sum(venta)*0.0125 from (
select SUM(vent) as venta from parley
union 
select SUM(vent) as venta from loteria
) as a


--
select fech, sistema.oid from archivo
join sistema on sistema.oid=archivo.tipo
where 
fech = '2017-02-01 00:00:00'
group by fech,sistema.oid
order by sistema.oid


--
select fech, count(archivo.oid) from archivo
group by fech
order by fech

