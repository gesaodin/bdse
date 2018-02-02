DROP TABLE IF EXISTS  archivo;
CREATE TABLE archivo
(
  oid serial NOT NULL,
  esta integer,
  nomb character varying(32) NOT NULL,
  fech timestamp without time zone,
  fcre timestamp without time zone,
  fpro timestamp without time zone,
  urls character varying(256),
  cant integer,
  resp integer,
  publ integer,
  tipo integer,
  tabl integer,
  acci character varying(32),
  CONSTRAINT archivo_pkey PRIMARY KEY (nomb)
);




DROP TABLE IF EXISTS figura;

CREATE TABLE figura
(
  oid serial NOT NULL,
  agen character varying(128),
  vent numeric,
  prem numeric,
  comi numeric,
  usua integer,
  fech timestamp without time zone,
  fcre timestamp without time zone,
  sist integer,
  arch integer,
  CONSTRAINT figura_pkey PRIMARY KEY (oid)
);


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
  sist integer,
  arch integer,
  CONSTRAINT loteria_pkey PRIMARY KEY (oid)
);

DROP TABLE IF EXISTS parley;

CREATE TABLE parley
(
  oid serial NOT NULL,
  agen character varying(128),
  vent numeric,
  prem numeric,
  comi numeric,
  usua integer,
  fech timestamp without time zone,
  fcre timestamp without time zone,
  sist integer,
  arch integer,
  CONSTRAINT parley_pkey PRIMARY KEY (oid)
);