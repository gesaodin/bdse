package balance

//Queda Establece la asignación de los clientes
type Queda struct {
	Fecha string
	Tipo  string
}

//GGlobal Caso 1:
//INCLUYENDO SALDOS ( + & -)
func (Q *Queda) GGlobal() string {
	return `SELECT 
	A.oid,obse,qued,freq,saldo
		FROM grupo A
		JOIN (
		SELECT grupo, SUM(saldo) saldo FROM (
		SELECT  zr.grupo, agencia.oid, agencia.obse, 
		SUM(vent)- SUM(prem) - SUM(comi) AS saldo, s.oid as programa, s.arch AS archivo 
		FROM (
			SELECT agen, fech, vent,prem,comi, sist from loteria
			UNION
			SELECT agen, fech, vent,prem,comi, sist from parley
			UNION
			SELECT agen, fech, vent,prem,comi, sist from figura
		) AS A
		JOIN zr_agencia zr ON A.agen=zr.codi
		JOIN agencia ON agencia.oid = zr.oida
		JOIN sistema s ON s.oid=A.sist
		WHERE A.fech ` + Q.Fecha + `
		GROUP BY zr.grupo, agencia.oid, agencia.obse, s.oid,  s.arch
		ORDER BY agencia.oid ) AS AGENCIA 
		GROUP BY grupo
		ORDER BY grupo ) --Fin Join inicial
		AS L ON A.oid=L.grupo
		-- AND qued > 0
		-- AND freq = ` + Q.Tipo
}

//GIndividual Caso 2:
//SOLO SALDOS ( + )
func (Q *Queda) GIndividual() string {
	return `SELECT 
		A.oid,obse,qued,freq,saldo
		FROM grupo A
		JOIN (
		SELECT grupo, SUM(saldo) saldo FROM (
		SELECT  zr.grupo, agencia.oid, agencia.obse, 
		SUM(vent)- SUM(prem) - SUM(comi) AS saldo
		FROM (
			SELECT agen, fech, vent,prem,comi, sist from loteria
			UNION
			SELECT agen, fech, vent,prem,comi, sist from parley
			UNION
			SELECT agen, fech, vent,prem,comi, sist from figura
		) AS A
		JOIN zr_agencia zr ON A.agen=zr.codi
		JOIN agencia ON agencia.oid = zr.oida
		JOIN sistema s ON s.oid=A.sist
		WHERE A.fech ` + Q.Fecha + `
		GROUP BY zr.grupo, agencia.oid, agencia.obse
		ORDER BY agencia.oid ) AS AGENCIA 
		WHERE saldo > 0
		GROUP BY grupo
		ORDER BY grupo ) --Fin Join inicial
		AS L ON A.oid=L.grupo
		WHERE saldo > 0 -- EN Agencia No por Grupo
		-- AND qued > 0
		-- AND freq = ` + Q.Tipo
}

//GPorPrograma Caso 3:
//SOLO SALDOS ( + )
func (Q *Queda) GPorPrograma() string {
	return `
	SELECT 
		A.oid,obse,qued,freq,saldo,L.sist
	FROM grupo A
	JOIN (
	
	SELECT  zr.grupo, A.sist,
	   SUM(vent)- SUM(prem) - SUM(comi) AS saldo
	FROM (
		SELECT agen, fech, vent,prem,comi, sist from loteria
		UNION
		SELECT agen, fech, vent,prem,comi, sist from parley
		UNION
		SELECT agen, fech, vent,prem,comi, sist from figura
	) AS A
	JOIN zr_agencia zr ON A.agen=zr.codi
	JOIN agencia ON agencia.oid = zr.oida
	JOIN sistema s ON s.oid=A.sist
	WHERE A.fech ` + Q.Fecha + `
	GROUP BY zr.grupo, A.sist
	ORDER BY zr.grupo
	) --Fin Join inicial
	AS L ON A.oid=L.grupo
	--WHERE saldo > 0 -- EN Agencia No por Grupo
	-- AND qued > 0
	-- AND freq = ` + Q.Tipo
}

//GPorJugada Caso 4:
//SOLO SALDOS ( + )
func (Q *Queda) GPorJugada() string {
	return `
		SELECT 
			L.grupo,qued,freq,saldo,L.arch
		FROM zr_negociacion_grupo_jugada A
		RIGHT JOIN (
		
		SELECT  zr.grupo, s.arch,
		SUM(vent)- SUM(prem) - SUM(comi) AS saldo
		FROM (
			SELECT agen, fech, vent,prem,comi, sist from loteria
			UNION
			SELECT agen, fech, vent,prem,comi, sist from parley
			UNION
			SELECT agen, fech, vent,prem,comi, sist from figura
		) AS A
		JOIN zr_agencia zr ON A.agen=zr.codi
		JOIN agencia ON agencia.oid = zr.oida
		JOIN sistema s ON s.oid=A.sist
		WHERE A.fech ` + Q.Fecha + `
		GROUP BY zr.grupo, s.arch
		ORDER BY zr.grupo
		) --Fin Join inicial
		AS L ON A.oid=L.grupo
		--WHERE saldo > 0 -- EN Agencia No por Grupo
		-- AND qued > 0
		-- AND freq = ` + Q.Tipo
}

/**************************************
**  Casos vistos desde la agencia
***************************************/

//AGlobal Caso 1:
//INCLUYENDO SALDOS ( + & -)
func (Q *Queda) AGlobal() string {
	return `SELECT 
		A.grupo,A.oid,obse,qued,freq,saldo
	FROM agencia A
	JOIN (
	SELECT grupo, oid, SUM(saldo) saldo FROM (
	SELECT  zr.grupo, agencia.oid, agencia.obse, 
	   SUM(vent)- SUM(prem) - SUM(comi) AS saldo, s.oid as programa, s.arch AS archivo 
	FROM (
		SELECT agen, fech, vent,prem,comi, sist from loteria
		UNION
		SELECT agen, fech, vent,prem,comi, sist from parley
		UNION
		SELECT agen, fech, vent,prem,comi, sist from figura
	) AS A
	JOIN zr_agencia zr ON A.agen=zr.codi
	JOIN agencia ON agencia.oid = zr.oida
	JOIN sistema s ON s.oid=A.sist
	WHERE A.fech  ` + Q.Fecha + `
	GROUP BY zr.grupo, agencia.oid, agencia.obse, s.oid,  s.arch
	ORDER BY agencia.oid ) AS AGENCIA 
	GROUP BY grupo, oid
	ORDER BY oid ) --Fin Join inicial
	AS L ON A.oid=L.oid
	WHERE saldo > 0 -- EN Agencia No por Grupo
	-- AND qued > 0
	-- AND freq = ` + Q.Tipo
}

//APorPrograma Caso 2:
//SOLO SALDOS ( + )
func (Q *Queda) APorPrograma() string {
	return `SELECT grupo, N.oid, N.obse, qued, freq, saldo,programa FROM zr_negociacion_agencia AS Z
	RIGHT JOIN (
	SELECT  zr.grupo, agencia.oid, agencia.obse, 
	   SUM(vent)- SUM(prem) - SUM(comi) AS saldo, s.oid as programa, s.arch AS archivo 
	FROM (
		SELECT agen, fech, vent,prem,comi, sist from loteria
		UNION
		SELECT agen, fech, vent,prem,comi, sist from parley
		UNION
		SELECT agen, fech, vent,prem,comi, sist from figura
	) AS A
	JOIN zr_agencia zr ON A.agen=zr.codi
	JOIN agencia ON agencia.oid = zr.oida
	JOIN sistema s ON s.oid=A.sist
	WHERE A.fech  ` + Q.Fecha + `
	GROUP BY zr.grupo, agencia.oid, agencia.obse, s.oid,  s.arch
	ORDER BY agencia.oid  ) AS N
	ON Z.oida=N.oid
	AND Z.oids=N.programa
	WHERE saldo > 0 -- EN Agencia No por Grupo
	-- AND qued > 0
	-- AND freq = ` + Q.Tipo + ` ORDER BY N.oid	`
}

//APorJugada Caso 3:
//SOLO SALDOS ( + )
func (Q *Queda) APorJugada() string {
	return `SELECT grupo, N.oid, N.obse, qued, freq, saldo,archivo 
	FROM zr_negociacion_agencia_jugada AS Z
	RIGHT JOIN (
	SELECT  zr.grupo, agencia.oid, agencia.obse, 
	   SUM(vent)- SUM(prem) - SUM(comi) AS saldo, s.arch AS archivo 
	FROM (
		SELECT agen, fech, vent,prem,comi, sist from loteria
		UNION
		SELECT agen, fech, vent,prem,comi, sist from parley
		UNION
		SELECT agen, fech, vent,prem,comi, sist from figura
	) AS A
	JOIN zr_agencia zr ON A.agen=zr.codi
	JOIN agencia ON agencia.oid = zr.oida
	JOIN sistema s ON s.oid=A.sist
	WHERE A.fech ` + Q.Fecha + `
	GROUP BY zr.grupo, agencia.oid, agencia.obse,  s.arch
	ORDER BY agencia.oid) AS N
	ON Z.oida=N.oid
	AND Z.oidt=N.archivo
	WHERE saldo > 0 -- EN Agencia No por Grupo
	-- AND qued > 0
	-- AND freq = ` + Q.Tipo + ` ORDER BY N.oid	`

}
