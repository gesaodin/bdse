package sys

import (
	"encoding/json"

	"../util"
)

var (
	Postgres bool
	Mysql    bool
	Mongodb  bool
)

const (
	ActivarConexionRemota      = true
	DesactivarConexionRemota   = false
	ActivarLogDeRegistro       = true
	DesactivarLogDeRegistro    = false
	ActivarRoles               = true
	DesactivarRoles            = false
	ActivarLimitEnConsultas    = true
	DesactivarLimitEnConsultas = false
	Puerto                     = 260804
	CodificacionDeArchivo      = "UTF-8"
	MaximoLimiteDeUsuarios     = 100
	MaximoLimiteDeConsultas    = 10
)

type BaseDeDatos struct {
	Modelo string
	Driver string
}

type CadenaDeConexion struct {
	Usuario   string
	BaseDatos string
	Clave     string
	Host      string
	Puerto    string
}

var Conexiones []CadenaDeConexion

func Iniciar() {
	var a util.Archivo
	a.NombreDelArchivo = "sys/config.json"
	data, _ := a.LeerTodo()
	e := json.Unmarshal(data, &Conexiones)
	util.Error(e)
}
