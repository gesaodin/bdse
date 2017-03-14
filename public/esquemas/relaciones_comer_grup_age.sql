
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
	CONSTRAINT agencia_obse_key UNIQUE (obse)
);
