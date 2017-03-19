
DROP TABLE IF EXISTS cobrosypagos_grupo;
CREATE TABLE cobrosypagos_grupo
(
  oid serial NOT NULL,
  oidg integer,
  fech timestamp without time zone,
  vien numeric,
  sald numeric,
  movi numeric,
  van numeric,
  CONSTRAINT cobrosypagos_grupo_pkey PRIMARY KEY (oid)
);
CREATE INDEX cobrosypagos_grupo_idgx ON cobrosypagos_grupo USING btree (oidg);
CREATE INDEX cobrosypagos_grupo_idfx ON cobrosypagos_grupo USING btree (fech);



DROP TABLE IF EXISTS cobrosypagoscierre_grupo;
CREATE TABLE cobrosypagoscierre_grupo
(
  oid serial NOT NULL,
  fech timestamp without time zone,
  esta integer,
  CONSTRAINT cobrosypagoscierre_grupo_pkey PRIMARY KEY (oid)
);