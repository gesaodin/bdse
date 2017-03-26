
delete from cobrosypagos where oid>4956
delete from cobrosypagos_grupo where oid>546

delete from cobrosypagos_agencia
insert into cobrosypagos_agencia  (oida,fech,vien,van)
select oida,'2016-12-31'::TIMESTAMP,vien,vien from cobrosypagos WHERE fech='2016-12-31'


insert into cobrosypagos_grupo  (oidg,fech,vien,van,movi,sald,erec)
SELECT grupo.oid,'2016-12-31'::TIMESTAMP, SUM(vien), SUM(vien),0,0,0 FROM grupo
JOIN agencia ON grupo.oid=agencia.grupo
JOIN cobrosypagos ON cobrosypagos.oida=agencia.oid
group by grupo.oid


select * from cobrosypagos_agencia  where oida =105