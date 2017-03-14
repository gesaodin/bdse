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
	comer integer,
	obse character varying(128),
	resp integer,
	fneg date,
	trip numeric,
	term numeric,
	qued numeric,
	part numeric,
	calc integer,
	freq integer,
	tipo integer, -- 1: GRUPO | 2: SUB GRUPO | 3 COLECTOR | 4 AGENCIA
	CONSTRAINT grupo_comer_fkey FOREIGN KEY (comer)
		REFERENCES comercializadora (oid) MATCH SIMPLE
		ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT grupo_obse_key UNIQUE (obse)
);

-- GRUPO, SUBGRUPO, COLECTOR, AGENCIA
DROP TABLE IF EXISTS zr_gsca_localizacion;
CREATE TABLE zr_gsca_localizacion 
(
	oid serial NOT NULL PRIMARY KEY,
	comer int,
	grupo int, -- CODIGO GRUPO, SUB, COLECTOR, AGENCIA
	subgr int,
	colec int,
	oida int NOT NULL,
	parro int,
	casa char varying(256),
	dire char varying(256),
	cuen char varying(23), -- 20 DIGITOS MAS EL FORMATO
	tele char varying(16),
	celu char varying(16),
	obse char varying(255),
	tipo int, -- 1: GRUPO | 2: SUB GRUPO | 3 COLECTOR | 4 AGENCIA
	fech timestamp without time zone
);
CREATE INDEX zr_gsca_localizacion_idxg ON zr_gsca_localizacion USING btree (grupo);
CREATE INDEX zr_gsca_localizacion_idxp ON zr_gsca_localizacion USING btree (parro);
-- delete from grupo where oid=637
-- delete from usuario where oid=2541
-- select * from usuario
--INSERT INTO grupo (comer, obse,tipo) VALUES (1,'ROAN',0),(1,'ROSITA',0),(1,'CINDY',0),(1,'COMPADRE',0),
--(1,'EIMAR',0),(1,'JOSE LH',0),(1,'JUNIOR',0), (1,'LEPE',0), (1,'MARYORI',0), (1,'ORLANDO',0),(1,'WILMER',0),
--(1,'YINNIS',0),(1,'ALEXPIER',0);


DROP TABLE IF EXISTS subgrupo;
CREATE TABLE subgrupo 
(
	oid serial NOT NULL PRIMARY KEY,	
	comer int, -- COMERCIALIZADORA
	grupo int, -- GRUPO
	obse character varying(128),
	fneg date,
	trip numeric,
	term numeric,
	qued numeric,
	part numeric,
	calc integer,
	freq integer,
	tipo integer, -- 1: GRUPO | 2: SUB GRUPO | 3 COLECTOR | 4 AGENCIA
	CONSTRAINT subgrupo_obse_key UNIQUE (obse)
);


DROP TABLE IF EXISTS agencia;
CREATE TABLE agencia 
(
	oid serial NOT NULL PRIMARY KEY,
	comer int, -- COMERCIALIZADORA
	grupo int, -- GRUPO
	subgr int, -- SUBGRUPO
	colec int, -- COLECTOR
	obse char varying(128),
	fneg date, -- FECHA DE NEGOCIACION
	trip numeric,
	term numeric,
	qued numeric,
	part numeric,
	calc integer,
	freq integer,
	tipo integer, -- 1: GRUPO | 2: SUB GRUPO | 3 COLECTOR | 4 AGENCIA
	CONSTRAINT subgrupo_obse_key UNIQUE (obse)
);

-- UPDATE agencia SET tipo=1,comer=1,grupo=0,subgr=0,colec=0

--INSERT INTO agencia (obse) VALUES ('APMEMMPPCD00400');
--INSERT INTO agencia (obse) VALUES ('APMEMMPPJR00100');
--INSERT INTO agencia (obse) VALUES ('APMEMMPPAP00500');
--INSERT INTO agencia (obse) VALUES ('APMEMMPPCD00100');
--INSERT INTO agencia (obse) VALUES ('APMEMMPPCD05500');

/**
* DESCRIBE LA RELACION ENTRE LA AGENCIA Y LAS TAQUILLAS
*/
DROP TABLE IF EXISTS zr_agencia;
CREATE TABLE zr_agencia 
(
	
	comer int,
	grupo int,
	subgr int,
	colec int,
	oida int NOT NULL,
	codi char varying(128), -- CAJA O TAQUILLA
	saldoactual numeric, --AL CIERRE
	CONSTRAINT zr_agencia_key UNIQUE (comer,grupo,subgr,colec,oida,codi)
);
CREATE INDEX zr_agencia_idx ON zr_agencia USING btree (comer,codi);


/**
INSERT INTO zr_agencia (comer,grupo,subgr,colec,oida,codi) VALUES 
(1,1,'APMEMMPPCD00400'), (1,1,'APMEMMPPCD00401');
INSERT INTO zr_agencia (comer,grupo,subgr,colec,oida,codi) VALUES 
(1,2,'APMEMMPPJR00100'), (1,2,'APMEMMPPJR00101');
INSERT INTO zr_agencia (comer,grupo,subgr,colec,oida,codi) VALUES 
(1,3,'APMEMMPPAP00500'), (1,3,'APMEMMPPAP00501'), (1,3,'APMEMMPPAP00502'), (1,3,'APMEMMPPAP00503'), (1,3,'APMEMMPPAP00504');
INSERT INTO zr_agencia (comer,grupo,subgr,colec,oida,codi) VALUES 
(1,3,'MAMEMMPPAP00500'), (1,3,'MAMEMMPPAP00501');
INSERT INTO zr_agencia (comer,grupo,subgr,colec,oida,codi) VALUES 
(1,3,'MMAMEMMPPAP00500'),(1,3,'MMAMEMMPPAP00501');
INSERT INTO zr_agencia (comer,grupo,subgr,colec,oida,codi) VALUES 
(1,4,'APMEMMPPCD00100');     
INSERT INTO zr_agencia (comer,grupo,subgr,colec,oida,codi) VALUES 
(1,5,'APMEMMPPCD05500'), (1,5,'MAMEMMPPCD05500'), (1,5,'MMAMEMMPPCD05500'),(1,5,'MAMEMMPPCD005500');
**/

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
	fech timestamp without time zone
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
	agen character varying(254), -- Agencia
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
	agen character varying(254), -- Agencia
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
	oida int,
	fech timestamp without time zone,
	vien numeric, --BALANCE DE REGISTRO
	sald numeric,
	movi numeric,
	van numeric, --BALANCE DE REGISTRO	
	CONSTRAINT cobrosypagos_pkey PRIMARY KEY (oid)
);
CREATE INDEX cobrosypagos_oida_idx ON cobrosypagos USING btree (oida);
CREATE INDEX cobrosypagos_fech_idx ON cobrosypagos USING btree (fech);


DROP TABLE IF EXISTS cobrosypagoscierre;
CREATE TABLE cobrosypagoscierre
(
	oid serial NOT NULL,
	fech timestamp without time zone,
	esta int,
	CONSTRAINT cobrosypagoscierre_pkey PRIMARY KEY (oid)
);
CREATE INDEX cobrosypagoscierre_esta_idx ON cobrosypagoscierre USING btree (esta);
CREATE INDEX cobrosypagoscierre_fech_idx ON cobrosypagoscierre USING btree (fech);


DROP TABLE IF EXISTS usuario;
CREATE TABLE usuario
(
  oid serial NOT NULL,
  nomb character varying(32),
  ncom character varying(255),
  corr character varying(255),
  fech timestamp without time zone,
  esta integer,
  rol character varying(255),
  toke character varying(255),
  CONSTRAINT usuario_pkey PRIMARY KEY (oid),
  CONSTRAINT usuario_pkey_nomb UNIQUE (nomb)
);

INSERT INTO usuario (nomb,ncom,corr,fech,esta,rol, toke) VALUES 
(
	'carlos', 'Administrador Del Sistema','carlos@admin.com',
	Now(), 1, 'Administrador', md5('carlosza63qj2p')
	
),
(
	'admin', 'Administrador Del Sistema','carlos@admin.com',
	Now(), 1, 'Administrador', md5('admin123')
	
);


DROP TABLE IF EXISTS sistema;
CREATE TABLE sistema (
	oid serial NOT NULL,
	obse char varying(256),
	arch integer,
	CONSTRAINT sistema_pkey PRIMARY KEY (oid)	
);
1
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
  comer int, -- Comercializadora
  grupo int, -- Grupo
  subgr int, -- Sub Grupo
  colec int, -- Colector
  oida int,
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
  resp character varying(254), --Respuesta
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
  comer int, -- Comercializadora
  grupo int, -- Grupo
  subgr int, -- Sub Grupo
  colec int, -- Colector
  oida int,
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
  resp character varying(254), --Respuesta
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
