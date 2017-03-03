DROP TABLE IF EXISTS persona;
CREATE TABLE persona 
(
	oid serial NOT NULL PRIMARY KEY,
	cedu char varying(32),
	nomb char varying(256) 	
);

DROP TABLE IF EXISTS comercializadora;
CREATE TABLE comercializadora 
(
	oid serial NOT NULL PRIMARY KEY,
	obse char varying(128),
	resp int,	
	telf char varying(16),
	saldoactual numeric --AL CIERRE
);
CREATE INDEX comercializadora_obse_idx ON comercializadora USING btree (obse COLLATE pg_catalog."default");
INSERT INTO comercializadora (obse) VALUES ('ALPIER');

DROP TABLE IF EXISTS grupo;
CREATE TABLE grupo 
(
	oid serial NOT NULL PRIMARY KEY,
	come int,
	obse char varying(128),
	resp int,	
	telf char varying(16),
	saldoactual numeric --AL CIERRE
);
CREATE INDEX grupo_obse_idx ON grupo USING btree (obse COLLATE pg_catalog."default");
INSERT INTO grupo (come, obse) VALUES (1,'ROAN'),(1,'ROSITA'),(1,'CINDY'),(1,'COMPADRE'),
(1,'EIMAR'),(1,'JOSE LH'),(1,'JUNIOR'), (1,'LEPE'), (1,'MARYORI'), (1,'ORLANDO'),(1,'WILMER'),
(1,'YINNIS'),(1,'ALEXPIER');


DROP TABLE IF EXISTS agencia;
CREATE TABLE agencia 
(
	oid serial NOT NULL PRIMARY KEY,
	-- comer int,
	-- grupo int,
	obse char varying(128),
	resp int,	
	telf char varying(16),
	saldoactual numeric --AL CIERRE
);
CREATE INDEX agencia_obse_idx ON agencia USING btree (obse COLLATE pg_catalog."default");

INSERT INTO agencia (obse) VALUES ('APMEMMPPCD00400');
INSERT INTO agencia (obse) VALUES ('APMEMMPPJR00100');
INSERT INTO agencia (obse) VALUES ('APMEMMPPAP00500');
INSERT INTO agencia (obse) VALUES ('APMEMMPPCD00100');
INSERT INTO agencia (obse) VALUES ('APMEMMPPCD05500');

/**
* DESCRIBE LA RELACION ENTRE LA AGENCIA Y LAS TAQUILLAS
*/
DROP TABLE IF EXISTS zr_agencia;
CREATE TABLE zr_agencia 
(
	oida int NOT NULL,
	codi char varying(128) 
);
CREATE INDEX zr_agencia_oida_idx ON zr_agencia USING btree (oida);
CREATE INDEX zr_agencia_codi_idx ON zr_agencia USING btree (codi COLLATE pg_catalog."default");

INSERT INTO zr_agencia (oida,codi) VALUES (1,'APMEMMPPCD00400'), (1,'APMEMMPPCD00401');
INSERT INTO zr_agencia (oida,codi) VALUES (2,'APMEMMPPJR00100'), (2,'APMEMMPPJR00101');
INSERT INTO zr_agencia (oida,codi) VALUES (3,'APMEMMPPAP00500'), (3,'APMEMMPPAP00501'), (3,'APMEMMPPAP00502'), (3,'APMEMMPPAP00503'), (3,'APMEMMPPAP00504');
INSERT INTO zr_agencia (oida,codi) VALUES (3,'MAMEMMPPAP00500'), (3,'MAMEMMPPAP00501');
INSERT INTO zr_agencia (oida,codi) VALUES (3,'MMAMEMMPPAP00500'),(3,'MMAMEMMPPAP00501');
INSERT INTO zr_agencia (oida,codi) VALUES (4,'APMEMMPPCD00100');     
INSERT INTO zr_agencia (oida,codi) VALUES (5,'APMEMMPPCD05500'), (5,'MAMEMMPPCD05500'), (5,'MMAMEMMPPCD05500'),(5,'MAMEMMPPCD005500');


DROP TABLE IF EXISTS zr_agencia_taquilla;
CREATE TABLE zr_agencia_taquilla(
	oid serial NOT NULL PRIMARY KEY,
	oida int,
	nomb char varying(128),
	fech timestamp without time zone
);


DROP TABLE IF EXISTS zr_agencia_sistema;
CREATE TABLE zr_agencia_sistema
(
	oida int,
	oids int,
	fech timestamp without time zone
	
);

DROP TABLE IF EXISTS zh_agencia_saldo;
CREATE TABLE zh_agencia_saldo 
(
	oid serial NOT NULL,
	codi char varying(128),
	fech timestamp without time zone,
	saldo numeric
	
);



INSERT INTO zh_agencia_saldo (codi, fech, saldo) VALUES ('APMEMMPPCD00400', now(), 20000);

DROP TABLE IF EXISTS loteria;
CREATE TABLE loteria
(
	oid serial NOT NULL,
	agen character varying(128),
	vent numeric,
	prem numeric,
	comi numeric,
	usua integer,
	fech timestamp without time zone,
	fcre timestamp without time zone,
	sist int,
	arch int,
	CONSTRAINT loteria_pkey PRIMARY KEY (oid)
);
CREATE INDEX loteria_fech_idx ON loteria USING btree (fech);
CREATE INDEX loteria_arch_idx ON loteria USING btree (arch);
CREATE INDEX loteria_sist_idx ON loteria USING btree (sist);


DROP TABLE IF EXISTS parley;
CREATE TABLE parley(
	oid serial NOT NULL,
	agen character varying(128),
	vent numeric,
	prem numeric,
	comi numeric,
	usua integer,
	fech timestamp without time zone,
	fcre timestamp without time zone,
	sist int,
	arch int,
	CONSTRAINT parley_pkey PRIMARY KEY (oid)
);
CREATE INDEX parley_fech_idx ON parley USING btree (fech);
CREATE INDEX parley_arch_idx ON parley USING btree (arch);
CREATE INDEX parley_sist_idx ON parley USING btree (sist);


DROP TABLE IF EXISTS cuentasxcobrar;
CREATE TABLE cuentasxcobrar
(
	oid serial NOT NULL,
	fech timestamp without time zone,
	cuen int,	
	mont numeric,
	sald numeric,
	
	CONSTRAINT cuentasxcobrar_pkey PRIMARY KEY (oid)
);
CREATE INDEX cuentasxcobrar_fech_idx ON cuentasxcobrar USING btree (fech);

DROP TABLE IF EXISTS banco;
CREATE TABLE banco
(
	oid serial NOT NULL PRIMARY KEY,
	nomb char varying(256), 
	nume char varying(4),
	auto char varying(256), -- Responsable o Autor
	cedu char varying(16),
	tipo int, -- Acepta Deposito (SI: Cuenta Bancaria NO : Cuenta Contable)
	fech timestamp without time zone,
);
INSERT INTO banco (nomb,auto,tipo,fech) VALUES 
('Venezuela',now(),'Alexander', 1),
('Provincial',now(),'Alexander', 1),
('Mercantil',now(),'Alexander', 1),
('Bicentenario',now(),'Alexander', 1),
('Banesco',now(),'Alexander', 1),
('BOD',now(),'Alexander', 1)


DROP TABLE IF EXISTS movimiento_ingreso;
CREATE TABLE movimiento_ingreso
(
	
	oid serial NOT NULL,
	comer int, -- Comercializadora
	grupo int, -- Grupo
	subgr int, -- Sub Grupo
	colec int, -- Colector
	agenc int, -- Agencia
	fech date,
	fapr timestamp without time zone,
	freg timestamp without time zone,
	tipo int, -- Tipo de Operacion	
	cuen int, -- Cuenta
	oper character varying(16),
	obse character varying(254),
	mont numeric,
	toke character varying(254),
	CONSTRAINT movimiento_ingreso_pkey PRIMARY KEY (oid)
);
CREATE INDEX movimiento_ingreso_toke_idx ON movimiento_ingreso USING btree (toke COLLATE pg_catalog."default");
CREATE INDEX movimiento_ingreso_fech_idx ON movimiento_ingreso USING btree (fech);
CREATE INDEX movimiento_ingreso_fapr_idx ON movimiento_ingreso USING btree (fapr);

DROP TABLE IF EXISTS movimiento_egreso;
CREATE TABLE movimiento_egreso
(
	oid serial NOT NULL,
	comer int, -- Comercializadora
	grupo int, -- Grupo
	subgr int, -- Sub Grupo
	colec int, -- Colector
	agenc int, -- Agencia
	fech date,
	fapr timestamp without time zone,
	freg timestamp without time zone,
	tipo int, -- Tipo de Operacion	
	cuen int, -- Cuenta
	oper character varying(16),
	obse character varying(254),
	mont numeric,
	toke character varying(254),
	CONSTRAINT movimiento_egreso_pkey PRIMARY KEY (oid)	
);
CREATE INDEX movimiento_egreso_toke_idx ON movimiento_egreso USING btree (toke COLLATE pg_catalog."default");
CREATE INDEX movimiento_egreso_fech_idx ON movimiento_egreso USING btree (fech);
CREATE INDEX movimiento_egreso_fapr_idx ON movimiento_egreso USING btree (fapr);

DROP TABLE IF EXISTS movimiento_prestamo;
CREATE TABLE movimiento_prestamo
(
	oid serial NOT NULL,
	agen char varying(256),
	tipo int,
	fech date,
	fapr timestamp without time zone,
	freg timestamp without time zone,
	mcuo int,	
	cuen int,
	saldo numeric,
	banc int,
	form int, -- Forma de Pago
	mont numeric,
	CONSTRAINT movimiento_prestamo_pkey PRIMARY KEY (oid)	
);
CREATE INDEX movimiento_prestamo_fech_idx ON movimiento_prestamo USING btree (fech);
CREATE INDEX movimiento_prestamo_fapr_idx ON movimiento_prestamo USING btree (fapr);


/**
* UN COBRO ESTA COMPUESTO POR MUCHOS MOVIMIENTOS
* LOS MOVIMIENTOS SON DE DIFERENTES TIPO 
*/
DROP TABLE IF EXISTS cobrosypagos;
CREATE TABLE cobrosypagos
(
	oid serial NOT NULL,
	agen char varying(256),
	fech timestamp without time zone,
	vien numeric, --BALANCE DE REGISTRO
	sald numeric,
	movi numeric,
	van numeric, --BALANCE DE REGISTRO	
	CONSTRAINT cobrosypagos_pkey PRIMARY KEY (oid)
);
CREATE INDEX cobrosypagos_fech_idx ON cobrosypagos USING btree (fech);



DROP TABLE IF EXISTS sistema;
CREATE TABLE sistema (
	oid serial NOT NULL,
	obse char varying(256),
	arch integer,
	CONSTRAINT sistema_pkey PRIMARY KEY (oid)	
);

DROP TABLE IF EXISTS usuario;
CREATE TABLE usuario
(
  oid serial NOT NULL,
  nomb character varying(32),
  ncom character varying(255),
  corr character varying(255),
  clav character varying(255),
  fech timestamp without time zone,
  esta integer,
  rol character varying(255),
  toke character varying(255),
  CONSTRAINT usuario_pkey PRIMARY KEY (oid)
);

INSERT INTO usuario (nomb,ncom,corr,clav,fech,esta,rol, toke) VALUES 
(
	'carlos', 'Administrador Del Sistema','carlos@admin.com',md5('za63qj2p'),
	Now(), 1, 'Administrador', md5('carlosza63qj2p')
	
),
(
	'admin', 'Administrador Del Sistema','carlos@admin.com',md5('za63qj2p'),
	Now(), 1, 'Administrador', md5('admin123')
	
),
(
	'APMEMMPPCD00400', 'APMEMMPPCD00400','carlos@admin.com',md5('123'),
	Now(), 1, 'Consumidor', md5('APMEMMPPCD00400123')
	
),
(
	'APMEMMPPAP00500', 'APMEMMPPAP00500','carlos@admin.com',md5('123'),
	Now(), 1, 'Consumidor', md5('APMEMMPPAP00500123')
	
);

DROP TABLE IF EXISTS archivo;
CREATE TABLE archivo
(
  oid serial NOT NULL,
  esta integer,
  nomb character varying(32) PRIMARY KEY,
  fech timestamp without time zone,
  fcre timestamp without time zone,
  fpro timestamp without time zone,
  urls character varying(256),
  cant integer,
  resp integer,
  publ integer,
  tipo integer,
  tabl integer,
  acci character varying(32) 
);

DROP TABLE IF EXISTS mensaje;
CREATE TABLE mensaje
(
  oid serial NOT NULL,
  orig character varying(32),
  dest character varying(254),
  msj character varying(254),
  fech timestamp without time zone,
  tipo integer,
  CONSTRAINT mensaje_pkey PRIMARY KEY (oid)
);



INSERT INTO sistema (obse,arch) VALUES ('Morpheus',0),('Pos1',0),('Pos2',0),
('Pos3',0),('Maticlo',0),('Ilbanquero',1),('CyberParley',1),
('Sport17',1),('Alien1',0), ('Alien2',0), ('Alien3',0), ('Turco1',0);


-- 0 DEBE | 1: HABER
DROP TABLE IF EXISTS debe;
CREATE TABLE debe
(
  oid serial NOT NULL,
  agen character varying(32),
  mont numeric,
  vouc character varying(254),
  fdep timestamp without time zone,
  freg timestamp without time zone,
  fope date,
  fapr date,
  tipo integer,
  banc integer,
  esta integer, -- 0 activos
  obse character varying(254),
  resp character varying(254), //Respuesta
  CONSTRAINT debe_pkey PRIMARY KEY (oid)
);
CREATE INDEX debe_fapr_idx ON debe USING btree (fapr);
CREATE INDEX debe_fest_idx ON debe USING btree (esta);
--INSERT INTO debe (agen,mont,vouc,fdep,freg,fope,tipo,banc,esta) VALUES ('APMEMMPPCD00400',2000,'VLO009888','2017-01-21',now(),'2017-01-19',1,1,0);


-- 0 DEBE | 1: HABER
DROP TABLE IF EXISTS haber;
CREATE TABLE haber
(
  oid serial NOT NULL,
  agen character varying(32),
  mont numeric,
  vouc character varying(254),
  fdep timestamp without time zone,--DEPOSITO
  freg timestamp without time zone,
  fope date,
  fapr date,
  tipo integer,
  banc integer,
  esta integer, -- 0 Pendiente Por Procesar  | 1 Activo
  obse character varying(254),
  resp character varying(254), //Respuesta
  CONSTRAINT haber_pkey PRIMARY KEY (oid)
);
CREATE INDEX haber_fapr_idx ON haber USING btree (fapr);
CREATE INDEX haber_fest_idx ON haber USING btree (esta);
/**
INSERT INTO haber (agen,mont,vouc,fdep,freg,fope,tipo,banc,esta) 
VALUES ('APMEBAPPRY00200',500,'VLO009888','2017-01-21',now(),'2017-01-21',1,1,0),
('APMEMMPPJR00100',7500,'VLO009888','2017-01-23',now(),'2017-01-23',1,1,0);
**/
DROP TABLE IF EXISTS saldos;
CREATE TABLE saldos
(
  oid serial NOT NULL PRIMARY KEY,
  obse character varying(128),
  resp integer,  
  saldoactual numeric
)


DROP TABLE IF EXISTS cuenta;
CREATE TABLE cuenta (
    cod int NOT NULL,
    nomb character varying(30),
    num character varying(20),
    tipo integer
);


